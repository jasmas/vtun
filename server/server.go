package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/tun"
	"github.com/patrickmn/go-cache"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

// Start server
func Start(config config.Config) {
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
	log.Printf("vtun server started on %v,CIDR is %v", config.LocalAddr, config.CIDR)
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
		srcAddr := srcAddr(b)
		dstAddr := dstAddr(b)
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
		dstAddr := dstAddr(b)
		srcAddr := srcAddr(b)
		if dstAddr == "" || srcAddr == "" {
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

func srcAddr(b []byte) string {
	if waterutil.IPv4Protocol(b) == waterutil.UDP || waterutil.IPv4Protocol(b) == waterutil.TCP {
		ip := waterutil.IPv4Source(b)
		port := waterutil.IPv4SourcePort(b)
		addr := fmt.Sprintf("%s:%d", ip.To4().String(), port)
		return addr
	} else if waterutil.IPv4Protocol(b) == waterutil.ICMP {
		ip := waterutil.IPv4Source(b)
		return ip.To4().String()
	}
	return ""
}

func dstAddr(b []byte) string {
	if waterutil.IPv4Protocol(b) == waterutil.UDP || waterutil.IPv4Protocol(b) == waterutil.TCP {
		ip := waterutil.IPv4Destination(b)
		port := waterutil.IPv4DestinationPort(b)
		addr := fmt.Sprintf("%s:%d", ip.To4().String(), port)
		return addr
	} else if waterutil.IPv4Protocol(b) == waterutil.ICMP {
		ip := waterutil.IPv4Destination(b)
		return ip.To4().String()
	}
	return ""
}
