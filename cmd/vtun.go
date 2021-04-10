package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/songgao/water"
)

const (
	BufferSize = 1500
	MTU        = "1300"
)

// New vtun
func New(local *string, remote *string, port *int, key *string) {
	os := runtime.GOOS
	if "linux" != os {
		log.Fatal("Only support linux!")
		return
	}
	hashKey := cipher.CreateHash(*key)
	// create tun
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = "vtun"
	iface, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	if nil != err {
		log.Fatalln("Unable to allocate TUN interface:", err)
	}
	log.Println("Interface allocated:", iface.Name())
	// config tun
	configTun(local, iface)
	// start udp listener
	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", *remote, *port))
	if nil != err {
		log.Fatalln("Unable to resolve remote addr:", err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", *port))
	if nil != err {
		log.Fatalln("Unable to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if nil != err {
		log.Fatalln("Unable to listen on UDP socket:", err)
	}
	defer conn.Close()
	// read data from remote and write data to tun
	go func() {
		buf := make([]byte, BufferSize)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil || n == 0 {
				break
			}
			b := buf[:n]
			// decrypt data
			cipher.Decrypt(&b, hashKey)
			iface.Write(b)
		}
	}()
	// read data from tun and write to remote
	packet := make([]byte, BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			break
		}
		b := packet[:n]
		// encrypt data
		cipher.Encrypt(&b, hashKey)
		conn.WriteToUDP(b, remoteAddr)
	}
}

func configTun(local *string, iface *water.Interface) {
	execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "mtu", MTU)
	execCmd("/sbin/ip", "addr", "add", *local, "dev", iface.Name())
	execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "up")
}

func execCmd(c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalln("Exec /sbin/ip error:", err)
	}
}
