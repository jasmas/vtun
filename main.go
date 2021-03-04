package main

import (
	"flag"
	"log"

	vtun "github.com/net-byte/vtun/cmd"
)

var (
	local  = flag.String("l", "172.16.0.1/24", "Local tun interface IP/MASK like 172.16.0.1/24")
	remote = flag.String("r", "172.16.0.2", "Remote server external IP like 172.16.0.2")
	port   = flag.Int("p", 2001, "UDP port")
	key    = flag.String("k", "6da62287-979a-4eb4-a5ab-8b3d89da134b", "Encrypt key")
)

func main() {
	flag.Parse()
	if "" == *local {
		flag.Usage()
		log.Fatalln("local ip is not specified")
	}
	if "" == *remote {
		flag.Usage()
		log.Fatalln("remote ip is not specified")
	}
	vtun.New(local, remote, port, key)
}
