package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func isMagicPacket(packet []byte, macAddr *string) bool {
	// fmt.Printf("%s\n", hex.EncodeToString(packet))
	if len(packet) != 102 {
		return false
	}
	m := hex.EncodeToString(packet[12:18])
	desired := strings.Repeat("ff", 6) + strings.Repeat(m, 16)
	if hex.EncodeToString(packet) == desired {
		*macAddr = strings.ToUpper(m[:2] + ":" + m[2:4] + ":" + m[4:6] + ":" + m[6:8] + ":" + m[8:10] + ":" + m[10:12])
		return true
	}
	return false
}

func sendPacket(addr string, port int, packet []byte) {
	raddr := net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	}
	conn, err := net.DialUDP("udp", nil, &raddr)
	if err != nil {
		fmt.Printf("Broadcast failed, err: %v\n", err)
	}
	defer conn.Close()
	conn.Write(packet)
}

func onRecvPacket(packet []byte) {
	var macAddr string
	if isMagicPacket(packet, &macAddr) {
		fmt.Printf("Magic %v ---> %s:%d ---> %s:%d (%s)\n", caddr, *addr, *port, *baddr, *bport, macAddr)
		sendPacket(*baddr, *bport, packet)
	}
}

func main() {
	defaultAddr := os.Getenv("WOL_ADDR")
	defaultPort, err := strconv.Atoi(os.Getenv("WOL_PORT"))
	if err != nil {
		defaultPort = 1999
	}
	defaultBAddr := os.Getenv("WOL_BADDR")
	defaultBPort, err := strconv.Atoi(os.Getenv("WOL_BPORT"))
	if err != nil {
		defaultBPort = 9
	}

	var addr = flag.String("addr", defaultAddr, "Listen address")
	var port = flag.Int("port", defaultPort, "Listen port")
	var baddr = flag.String("baddr", defaultBAddr, "Broadcast address")
	var bport = flag.Int("bport", defaultBPort, "Broadcast port")

	flag.Parse()
	listen, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", *addr, *port))
	if err != nil {
		fmt.Printf("Listen failed, err: %v\n", err)
		return
	}
	fmt.Printf("Listen %s:%d ---> %s:%d\n", *addr, *port, *baddr, *bport)
	defer listen.Close()
	for {
		var data [4096]byte
		n, caddr, err := listen.ReadFrom(data)
		if err != nil {
			fmt.Printf("Read udp failed, err: %v\n", err)
			continue
		}
		go onRecvPacket(data[:n])
	}
}
