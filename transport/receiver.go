package transport

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

func Receive() {
	addr, err := net.ResolveUDPAddr("udp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(buf[:n]))
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, uint32(1))
		conn.WriteToUDP(buf.Bytes(), clientAddr)
	}
}
