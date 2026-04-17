package transport

import (
	"log"
	"net"
	"time"

	"github.com/jsndz/rldp/protocol"
	"github.com/jsndz/rldp/types"
)

func Receive() {
	addr, err := net.ResolveUDPAddr("udp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	buffer := make(map[uint32]types.Frame)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	var counter uint32

	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		seq, ack, _, payload, _ := protocol.Decoding(buf[:n])
		if !protocol.VerifyChecksum(buf[:n]) {
			log.Println("received corrupted frame with seq:", seq, "ack:", ack, "payload:", payload)
			continue
		}
		if seq < (counter + 1) {
			log.Println("received duplicate frame with seq:", seq, "expected:", counter+1)
			resp := protocol.Encoding(types.Frame{
				Seq:  uint32(seq),
				Ack:  uint32(seq),
				Type: uint8(types.ACK),
			})
			conn.WriteToUDP(resp, clientAddr)
			continue
		}
		if seq > (counter + 1) {
			log.Println("received out of order frame with seq:", seq, "expected:", counter+1)

			buffer[seq] = types.Frame{
				Seq:     uint32(seq),
				Ack:     uint32(seq),
				Type:    uint8(types.DATA),
				Payload: []byte(payload),
			}
			resp := protocol.Encoding(types.Frame{
				Seq:  uint32(seq),
				Ack:  uint32(seq),
				Type: uint8(types.ACK),
			})
			conn.WriteToUDP(resp, clientAddr)
			continue
		}
		log.Println("received data:", payload, "from", clientAddr)
		resp := protocol.Encoding(types.Frame{
			Seq:  uint32(seq),
			Ack:  uint32(seq),
			Type: uint8(types.ACK),
		})
		conn.WriteToUDP(resp, clientAddr)
		counter++
		for {
			frame, ok := buffer[(counter)]
			if !ok {
				break
			}
			delete(buffer, (counter))
			resp := protocol.Encoding(types.Frame{
				Seq:  uint32(frame.Seq),
				Ack:  uint32(frame.Seq),
				Type: uint8(types.ACK),
			})
			time.Sleep(time.Second * 9)
			conn.WriteToUDP(resp, clientAddr)
			log.Println("processed buffered frame with seq:", frame.Seq, "payload:", string(frame.Payload))
			counter++
		}
	}
}
