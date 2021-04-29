package netutil

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/net-byte/vtun/common/config"
	"github.com/songgao/water/waterutil"
)

func SrcAddr(b []byte) (addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			addr = ""
		}
	}()
	if waterutil.IPv4Protocol(b) == waterutil.UDP || waterutil.IPv4Protocol(b) == waterutil.TCP {
		ip := waterutil.IPv4Source(b)
		port := waterutil.IPv4SourcePort(b)
		addr = fmt.Sprintf("%s:%d", ip.To4().String(), port)
		log.Printf("SrcAddr %v", addr)
		return addr
	} else if waterutil.IPv4Protocol(b) == waterutil.ICMP {
		ip := waterutil.IPv4Source(b)
		return ip.To4().String()
	}
	return addr
}

func DstAddr(b []byte) (addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			addr = ""
		}
	}()
	if waterutil.IPv4Protocol(b) == waterutil.UDP || waterutil.IPv4Protocol(b) == waterutil.TCP {
		ip := waterutil.IPv4Destination(b)
		port := waterutil.IPv4DestinationPort(b)
		addr = fmt.Sprintf("%s:%d", ip.To4().String(), port)
		log.Printf("DstAddr %v", addr)
		return addr
	} else if waterutil.IPv4Protocol(b) == waterutil.ICMP {
		ip := waterutil.IPv4Destination(b)
		return ip.To4().String()
	}
	return addr
}

func ConnectWS(config config.Config) *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: config.ServerAddr, Path: "/way-to-freedom"}
	header := make(http.Header)
	header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("[client] failed to dial websocket %v", err)
		return nil
	}
	return c
}

func CloseWS(wsConn *websocket.Conn) {
	wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*5))
	wsConn.Close()
}
