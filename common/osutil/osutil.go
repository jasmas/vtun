package osutil

import (
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func ConfigTun(cidr string, iface *water.Interface) {
	execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "mtu", "1300")
	execCmd("/sbin/ip", "addr", "add", cidr, "dev", iface.Name())
	execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "up")
}

func execCmd(c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalln("failed to exec /sbin/ip error:", err)
	}
}
