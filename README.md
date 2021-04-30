# vtun

A simple vpn.  

[![Travis](https://travis-ci.com/net-byte/vtun.svg?branch=master)](https://github.com/net-byte/vtun)
[![Go Report Card](https://goreportcard.com/badge/github.com/net-byte/vtun)](https://goreportcard.com/report/github.com/net-byte/vtun)
![image](https://img.shields.io/badge/License-MIT-orange)
![image](https://img.shields.io/badge/License-Anti--996-red)

# Usage  

```
Usage of ./vtun:
  -S    server mode
  -c string
        tun interface CIDR (default "172.16.0.1/24")
  -k string
        encryption key (default "6w9z$C&F)J@NcRfWjXn3r4u7x!A%D*G-")
  -p string
        protocol udp/ws (default "udp")
  -l string
        local address (default "0.0.0.0:3000")
  -s string
        server address (default "0.0.0.0:3001")
```  

## Client

```
sudo ./vtun -l=:3000 -s=server-addr:3001 -c=172.16.0.10/24 -k=123456

```
## Server

```
sudo ./vtun -S -l=:3001 -c=172.16.0.1/24 -k=123456

```

## Enable IP forwarding on server

```
sudo echo 1 > /proc/sys/net/ipv4/ip_forward
sudo iptables -t nat -A POSTROUTING -s 172.16.0.0/24 -o ens3 -j MASQUERADE
sudo sysctl -p
```
