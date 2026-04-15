package transport

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

const (
	DATA = 1
	ACK  = 2
)

func Send(data string, address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	// resolves the address into a format with ip and port
	if err != nil {
		log.Fatal(err)
	}
	seq := uint32(1)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, uint32(seq))
	binary.Write(&buf, binary.BigEndian, uint32(0))
	buf.WriteByte(byte(DATA))
	buf.WriteString(data)

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("sent data to", address)
	var responseBuf bytes.Buffer
	for {

		_, _, err := conn.ReadFromUDP(responseBuf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("received response from", address)
	}
}
