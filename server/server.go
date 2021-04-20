package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/tun"
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
	forwarder := &Forwarder{localConn: conn}
	go forwarder.forward(iface, conn)
	// read data from client
	buf := make([]byte, 1500)
	for {
		n, cliAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			forwarder.cliConnMap.Delete(cliAddr.String())
			continue
		}
		b := buf[:n]
		// decrypt data
		cipher.Decrypt(&b)
		iface.Write(b)
		forwarder.cliConnMap.LoadOrStore(cliAddr.String(), cliAddr)
	}
}

type Forwarder struct {
	localConn  *net.UDPConn
	cliConnMap sync.Map
}

func (f *Forwarder) forward(iface *water.Interface, conn *net.UDPConn) {
	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		fd := ForwardData{localConn: conn, data: packet[:n]}
		f.cliConnMap.Range(fd.walk)
	}
}

type ForwardData struct {
	localConn *net.UDPConn
	data      []byte
}

func (f *ForwardData) walk(key, value interface{}) bool {
	if waterutil.IsIPv4(f.data) {
		ip := waterutil.IPv4Destination(f.data)
		port := waterutil.IPv4DestinationPort(f.data)
		cliAddr := fmt.Sprintf("%s:%d", ip.To4().String(), port)
		log.Printf("to client:%v", cliAddr)
		if cliAddr == key.(string) {
			// encrypt data
			cipher.Encrypt(&f.data)
			f.localConn.WriteToUDP(f.data, value.(*net.UDPAddr))
			return false
		}
	} else {
		// encrypt data
		cipher.Encrypt(&f.data)
		f.localConn.WriteToUDP(f.data, value.(*net.UDPAddr))
	}
	return true
}
