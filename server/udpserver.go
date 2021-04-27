package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/common/netutil"
	"github.com/net-byte/vtun/tun"
	"github.com/patrickmn/go-cache"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

// StartUDPServer start udp server
func StartUDPServer(config config.Config) {
	if r := recover(); r != nil {
		log.Println("recover", r)
	}
	config.Init()
	iface := tun.CreateTun(config.CIDR)
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("vtun udp server started on %v,CIDR is %v", config.LocalAddr, config.CIDR)
	// forward data to client
	forwarder := &Forwarder{localConn: conn, connCache: cache.New(30*time.Minute, 10*time.Minute)}
	go forwarder.forward(iface, conn)
	// read data from client
	buf := make([]byte, 1500)
	for {
		n, cliAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			continue
		}
		// decrypt data
		b := cipher.Decrypt(buf[:n])
		if !waterutil.IsIPv4(b) {
			continue
		}
		iface.Write(b)
		srcAddr := netutil.SrcAddr(b)
		dstAddr := netutil.DstAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		forwarder.connCache.Set(key, cliAddr, cache.DefaultExpiration)
	}
}

type Forwarder struct {
	localConn *net.UDPConn
	connCache *cache.Cache
}

func (f *Forwarder) forward(iface *water.Interface, conn *net.UDPConn) {
	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		b := packet[:n]
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr := netutil.SrcAddr(b)
		dstAddr := netutil.DstAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		v, ok := f.connCache.Get(key)
		if ok {
			// encrypt data
			b = cipher.Encrypt(b)
			f.localConn.WriteToUDP(b, v.(*net.UDPAddr))
		}
	}
}
