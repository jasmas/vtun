package main

import (
	"flag"

	"github.com/net-byte/vtun/client"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/server"
)

func main() {
	config := config.Config{}
	flag.StringVar(&config.CIDR, "c", "172.16.0.0/24", "tun interface CIDR")
	flag.StringVar(&config.LocalAddr, "l", "0.0.0.0:3000", "local address")
	flag.StringVar(&config.ServerAddr, "s", "0.0.0.0:3001", "server address")
	flag.StringVar(&config.Key, "k", "6w9z$C&F)J@NcRfWjXn3r4u7x!A%D*G-", "encryption key")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")
	flag.Parse()

	if config.ServerMode {
		server.Start(config)
	} else {
		client.Start(config)
	}

}
