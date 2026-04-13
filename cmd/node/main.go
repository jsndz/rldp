package main

import (
	"os"
	"strings"

	"github.com/jsndz/rldp/transport"
)

func main() {
	args := os.Args

	switch args[1] {
	case "send":
		transport.Send(strings.Join(args[3:], " "), args[2])
	}

	go func() {
		transport.Receive()
	}()
}
