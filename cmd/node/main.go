package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jsndz/rldp/transport"
)

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
		transport.Send(strings.Join(args[3:], " "), args[2])
	case "receive":

		transport.Receive()

	default:
		fmt.Println("unknown command:", args[1])
		fmt.Println("usage:")
		fmt.Println("  go run cmd/node/main.go receive")
		fmt.Println("  go run cmd/node/main.go send <dst-mac> <message>")
	}
}
