package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/tun"
	"github.com/patrickmn/go-cache"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1500,
	WriteBufferSize:   1500,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StartWSServer start ws server
func StartWSServer(config config.Config) {
	config.Init()
	iface := tun.CreateTun(config.CIDR)
	c := cache.New(30*time.Minute, 10*time.Minute)
	go tunToWs(iface, c)
	log.Printf("vtun ws server started on %v,CIDR is %v", config.LocalAddr, config.CIDR)
	http.HandleFunc("/way-to-freedom", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		wsToTun(wsConn, iface, c)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello，世界！")
	})

	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(req.RemoteAddr, ":")[0]
		}
		resp := fmt.Sprintf("%v", ip)
		io.WriteString(w, resp)
	})

	http.ListenAndServe(config.ServerAddr, nil)
}

func closeWS(wsConn *websocket.Conn) {
	wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*5))
	wsConn.Close()
}

func tunToWs(iface *water.Interface, c *cache.Cache) {
	buffer := make([]byte, 1500)
	for {
		n, err := iface.Read(buffer)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		b := buffer[:n]
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr := srcAddr(b)
		dstAddr := dstAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		v, ok := c.Get(key)
		if ok {
			b = cipher.Encrypt(b)
			v.(*websocket.Conn).WriteMessage(websocket.BinaryMessage, b)
		}
	}
}

func wsToTun(wsConn *websocket.Conn, iface *water.Interface, c *cache.Cache) {
	defer closeWS(wsConn)
	for {
		wsConn.SetReadDeadline(time.Now().Add(time.Duration(30) * time.Second))
		_, b, err := wsConn.ReadMessage()
		if err != nil || err == io.EOF {
			break
		}
		b = cipher.Decrypt(b)
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr := srcAddr(b)
		dstAddr := dstAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		c.Set(key, wsConn, cache.DefaultExpiration)
		iface.Write(b[:])
	}
}
