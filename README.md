# go-tun2socks

[![Build Status](https://travis-ci.com/eycorsican/go-tun2socks.svg?branch=master)](https://travis-ci.com/eycorsican/go-tun2socks)

A tun2socks implementation written in Go.

Tested and worked on macOS, Linux, Windows and iOS (as a library).

## Overview

```
                                      core.NewLWIPStack()
                                           +
                                           |
                                           |
                                           |
                                           |                TCP/UDP             core.RegisterTCPConnectionHandler()
                                           |
                          core.Input()     |           core.Connection          core.RegisterUDPConnectionHandler()
                                           v
Application +------> TUN +-----------> lwIP stack +------------------------------> core.ConnectionHandler +-------> Remote proxy server +--> Destination


                         <-----------+
                    core.RegisterOutputFn()

```

## Features

- Support both TCP and UDP
- Support both IPv4 and IPv6
- Support ICMP local echoing with configurable packet delay time
- Support proxy handlers: `SOCKS5`, `Shadowsocks`, `V2Ray` (DNS cache is enabled in these handlers by default)
- Dynamically adding routing rules according to V2Ray's routing results (V2Ray proxy handler only)

## Build

`go-tun2socks` is using `cgo`, thus a C compiler is required.

```sh
go get github.com/eycorsican/go-tun2socks
cd $GOPATH/src/github.com/eycorsican/go-tun2socks
go get -d ./...
make clean && make build
./build/tun2socks -h
```

An alternative way to build (or cross compile) tun2socks is to use [`xgo`](https://github.com/karalabe/xgo), to use `xgo`, you also need `docker`:

```sh
# install docker: https://docs.docker.com/install

# install xgo
go get github.com/karalabe/xgo

go get github.com/eycorsican/go-tun2socks
cd $GOPATH/src/github.com/eycorsican/go-tun2socks
go get -d ./...
make clean && make xbuild
ls ./build
```

## Run

```sh
./build/tun2socks -tunName tun1 -tunAddr 240.0.0.2 -tunGw 240.0.0.1 -proxyType socks -proxyServer 1.2.3.4:1086
```

Note that the TUN device may have a different name, and it should be a different name on Windows unless you have renamed it, so make sure use `ifconfig`, `ipconfig` or `ip addr` to check it out.

## Create TUN device and Configure Routing Table

Suppose your original gateway is 192.168.0.1. The proxy server address is 1.2.3.4.

The following commands will need root permissions.

### macOS

The program will automatically create a TUN device for you on macOS. To show the created TUN device, use ifconfig.

Delete original gateway:

```sh
route delete default
```

Add our TUN interface as the default gateway:

```sh
route add default 240.0.0.1
```

Add a route for your proxy server to bypass the TUN interface:

```sh
route add 1.2.3.4/32 192.168.0.1
```

### Linux

The program will not create the TUN device for you on Linux. You need to create the TUN device by yourself:

```sh
ip tuntap add mode tun dev tun1
ip addr add 240.0.0.1 dev tun1
ip link set dev tun1 up
```

Delete original gateway:

```sh
ip route del default
```

Add our TUN interface as the default gateway:

```sh
ip route add default via 240.0.0.1
```

Add a route for your proxy server to bypass the TUN interface:

```sh
ip route add 1.2.3.4/32 via 192.168.0.1
```

### Windows

To create a TUN device on Windows, you need [Tap-windows](https://openvpn.net/index.php/download/community-downloads.html), refer [here](https://code.google.com/archive/p/badvpn/wikis/tun2socks.wiki) for more information.

Add our TUN interface as the default gateway:

```sh
# Using 240.0.0.1 is not allowed on Windows, we use 10.0.0.1 instead
route add 0.0.0.0 mask 0.0.0.0 10.0.0.1 metric 6
```

Add a route for your proxy server to bypass the TUN interface:

```sh
route add 1.2.3.4 192.168.0.1 metric 5
```

## A few notes for using V2Ray proxy handler
- Using V2Ray proxy handler: `tun2socks -proxyType v2ray -vconfig config.json`
- V2Ray proxy handler dials connections with a [V2Ray Instance](https://github.com/v2ray/v2ray-core/blob/master/functions.go)
- Configuration file V2Ray must in JSON format
- Proxy server addresses in the configuration file should be IPs and not domains except your system DNS will match "direct" rules
- Configuration file should not contain direct `domain` rules, since they cause infinitely looping requests
- Dynamic routing happens prior to packets input to lwIP, the [V2Ray Router](https://github.com/v2ray/v2ray-core/blob/master/features/routing/router.go) is used to check if the IP packet matching "direct" tag, information available for the matching process are (protocol, destination ip, destination port)
- To enable dynamic routing, just set the `-gateway` argument, for example: `tun2socks -proxyType v2ray -vconfig config.json -gateway 192.168.0.1`
- The tag "direct" is hard coded to identify direct rules, which if dynamic routing is enabled, will indicate adding routes to the original gateway for the corresponding IP packets
- Inbounds are not necessary

## TODO
- Built-in routing rules and routing table management
- Support ICMP packets forwarding

## This project is using lwIP 

This project is using a modified version of lwIP, you can checkout this repo to find out what are the changes: https://github.com/eycorsican/lwip

## Many thanks to the following projects
- https://savannah.nongnu.org/projects/lwip
- https://github.com/ambrop72/badvpn
- https://github.com/zhuhaow/tun2socks
- https://github.com/yinghuocho/gotun2socks
- https://github.com/shadowsocks/go-shadowsocks2
- https://github.com/nadoo/glider
