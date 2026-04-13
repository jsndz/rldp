package transport

import (
	"log"
	"net"

	"golang.org/x/sys/unix"
)

func Receive() {
	println("receiving data")
	iface, _ := net.InterfaceByName("wlan0")

	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(unix.ETH_P_ALL))
	// fd indicates raw packet socket
	// it is what is used to identify your socket in os
	if err == nil {
		log.Fatal(err)
	}
	// addr of the wlan0 in the format that fd can take
	addr := &unix.SockaddrLinklayer{
		Protocol: unix.ETH_P_ALL,
		Ifindex:  iface.Index,
	}

	if err := unix.Bind(fd, addr); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 65536)
	for {
		_, _, err := unix.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(buf[14:]))
	}
}
