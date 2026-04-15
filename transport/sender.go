package transport

import (
	"log"
	"net"
)

func Send(data string, address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	// resolves the address into a format with ip and port
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("sent data to", address)
}
