package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jsndz/rldp/protocol"
	"github.com/jsndz/rldp/transport"
	"github.com/jsndz/rldp/types"
)

const windowSize = 5
const timeout = 3 * time.Second

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  go run cmd/node/main.go receive")
		fmt.Println("  go run cmd/node/main.go send <dst-mac> <message>")
		return
	}

	switch args[1] {
	case "send":
		if len(args) < 4 {
			fmt.Println("usage: go run cmd/node/main.go send <dst-mac> <message>")
			return
		}
		addr := args[2]

		counter := 1
		for {
			var msg string
			fmt.Println(">")
			fmt.Scanln(&msg)

			if msg == "exit" {
				break
			}

			transport.Send(msg, addr, counter)
			counter++
		}
	case "receive":

		transport.Receive()
	case "batch-send":
		if len(args) < 4 {
			fmt.Println("usage: go run cmd/node/main.go batch-send <dst-mac> <message>")
			return
		}
		address := args[2]
		messages := args[3:]
		addr, err := net.ResolveUDPAddr("udp", address)
		// resolves the address into a format with ip and port
		if err != nil {
			log.Fatal(err)
		}
		var mu sync.Mutex
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		ackChan := make(chan uint32)
		buffer := make(map[int][]byte)
		counter := 1
		timer := time.NewTimer(time.Second * 3)
		go func() {
			recvBuf := make([]byte, 2048)
			for {
				n, _, err := conn.ReadFromUDP(recvBuf)
				if err != nil {
					log.Print("error", err)
				}
				seq, ack, frameType, payload, _ := protocol.Decoding(recvBuf[:n])
				if !protocol.VerifyChecksum(recvBuf[:n]) {
					log.Println("received corrupted frame with seq:", seq, "ack:", ack, "payload:", payload)
				}
				if frameType == types.ACK && ack == uint32(counter) {
					ackChan <- ack
				} else {
					log.Println("received non-ACK response from", ack, "with payload:", payload)

				}
			}
		}()
		base := (1)
		nextSeq := 1
		for base <= len(messages) {
			for nextSeq < windowSize+base && nextSeq <= len(messages) {
				buf := protocol.Encoding(types.Frame{
					Seq:     uint32(nextSeq),
					Ack:     0,
					Type:    types.DATA,
					Payload: []byte(messages[nextSeq-1]),
				})
				conn.Write(buf)
				mu.Lock()
				buffer[nextSeq] = buf
				mu.Unlock()
				if base == nextSeq {
					timer.Reset(timeout)
				}
				nextSeq++
			}
			select {
			case ack := <-ackChan:
				if int(ack) >= base {
					base = int(ack) + 1
					timer.Reset(timeout)
				}
			case <-timer.C:
				mu.Lock()
				for i := base; i < nextSeq; i++ {
					conn.Write(buffer[i])
					log.Println("resent:", i)
				}
				mu.Unlock()

				timer.Reset(timeout)
			}
		}
	default:
		fmt.Println("unknown command:", args[1])
		fmt.Println("usage:")
		fmt.Println("  go run cmd/node/main.go receive")
		fmt.Println("  go run cmd/node/main.go send <dst-mac> <message>")
	}
}
