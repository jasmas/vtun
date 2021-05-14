package register

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var _register *cache.Cache

func init() {
	_register = cache.New(30*time.Minute, 3*time.Minute)
}

func AddClientIP(ip string) {
	_register.Add(ip, 0, cache.DefaultExpiration)
}

func DeleteClientIP(ip string) {
	_register.Delete(ip)
}

func ExistClientIP(ip string) bool {
	_, ok := _register.Get(ip)
	return ok
}

func KeepAliveClientIP(ip string) {
	if ExistClientIP(ip) {
		_register.Increment(ip, 1)
	} else {
		AddClientIP(ip)
	}
}

func PickClientIP(cidr string) (clientIP string, prefixLength string) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}
	AddClientIP(ip.To4().String())
	pickIP := ipNet.IP.To4()
	for {
		pickIP[3]++
		if pickIP[3] >= 255 {
			break
		}
		if !ExistClientIP(pickIP.String()) {
			AddClientIP(pickIP.String())
			return pickIP.String(), strings.Split(cidr, "/")[1]
		}
	}
	return "", ""
}

func ListClientIP() []string {
	result := []string{}
	for k := range _register.Items() {
		result = append(result, k)
	}
	return result
}
