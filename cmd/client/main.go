package main

import (
	"log"
	"net"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

func main() {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = "vnet_demo"

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	var frame ethernet.Frame

	server, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Fatal(err)
		}

		nWrite, err := conn.Write([]byte(frame))
		if err != nil {
			log.Fatal(err)
		}

		log.Println("sent frame to server with length", nWrite)

		frame = frame[:n]
		log.Printf("Dst: %s\n", frame.Destination())
		log.Printf("Src: %s\n", frame.Source())
		log.Printf("Ethertype: % x\n", frame.Ethertype())
		log.Printf("Payload: % x\n", frame.Payload())
	}
}
