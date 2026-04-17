package transport

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jsndz/rldp/protocol"
	"github.com/jsndz/rldp/types"
)

func Send(data string, address string, counter int) (int, int) {
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

	buf := protocol.Encoding(types.Frame{
		Seq:     uint32(counter),
		Ack:     uint32(0),
		Type:    uint8(types.DATA),
		Payload: []byte(data),
	})

	_, err = conn.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("sent data to", address)

	for {
		// adding deadline if no response is available
		conn.SetDeadline(time.Now().Add(time.Second * 10))
		recvBuf := make([]byte, 2048)

		n, _, err := conn.ReadFromUDP(recvBuf)
		if err != nil {
			fmt.Println("sending again")
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				_, err = conn.Write(buf)
				if err != nil {
					return counter, 0
				}
				continue
			}
		}
		seq, ack, frameType, payload, _ := protocol.Decoding(recvBuf[:n])
		if !protocol.VerifyChecksum(recvBuf[:n]) {
			log.Println("received corrupted frame with seq:", seq, "ack:", ack, "payload:", payload)
			continue
		}
		if frameType == types.ACK && ack == uint32(counter) {
			log.Println("received ACK from", seq, "with payload:", payload)
		} else {
			log.Println("received non-ACK response from", ack, "with payload:", payload)
			continue
		}
		return counter, int(ack)

	}

}
