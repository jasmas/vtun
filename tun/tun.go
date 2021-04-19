package tun

import (
	"log"
	"runtime"

	"github.com/net-byte/vtun/common/osutil"
	"github.com/songgao/water"
)

func CreateTun(cidr string) (iface *water.Interface) {
	c := water.Config{}
	os := runtime.GOOS
	if os == "darwin" {
		c.DeviceType = water.TUN
		c.Name = "utun000"
	} else {
		c.DeviceType = water.TAP
		c.Name = "vtun"
	}
	iface, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to allocate TUN interface:", err)
	}
	log.Println("interface allocated:", iface.Name())
	osutil.ConfigTun(cidr, iface)
	return iface
}
