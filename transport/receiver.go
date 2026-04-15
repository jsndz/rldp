package transport

import (
	"log"
	"net"

	"github.com/jsndz/rldp/protocol"
	"github.com/jsndz/rldp/types"
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
		seq, _, _, payload := protocol.Decoding(buf[:n])
		resp := protocol.Encoding(types.Frame{
			Seq:     uint32(seq),
			Ack:     uint32(1),
			Type:    uint8(types.ACK),
			Payload: []byte("ACK for " + payload),
		})
		conn.WriteToUDP(resp, clientAddr)
	}

}
