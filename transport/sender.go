package transport

import (
	"encoding/binary"
	"log"
	"net"

	"golang.org/x/sys/unix"
)

func Send(data string, mac string) {

	// get interface of wifi wlan0 as the default
	iface, err := net.InterfaceByName("wlan0")
	if err != nil {
		panic(err)
	}

	// Binds directly to the network interface level
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(unix.ETH_P_ALL))
	if err != nil {
		log.Fatal(err)
	}
	defer unix.Close(fd)
	// get destination and src mac address
	dstMAC, _ := net.ParseMAC(mac) //-> 6 bytes
	// this is from the interface
	srcMac := iface.HardwareAddr // 6 -> bytes

	etherType := uint16(0x88B5) // custom ether type for rldp -> 2 bit
	// ethertype identifies which protocol is used to encapsulate the payload
	payload := []byte(data)

	frame := make([]byte, 14+len(payload)) // Ethernet frame is 14 bytes header + payload

	copy(frame[0:6], dstMAC)
	copy(frame[6:12], srcMac)
	binary.BigEndian.PutUint16(frame[12:14], etherType) // cant put directly the 16 bytes so
	copy(frame[14:], payload)

	// need address to be of this type to send
	addr := &unix.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Halen:    6,
		Addr:     [8]byte{dstMAC[0], dstMAC[1], dstMAC[2], dstMAC[3], dstMAC[4], dstMAC[5]},
		Protocol: unix.ETH_P_ALL,
	}

	err = unix.Sendto(fd, frame, 0, addr)
	if err != nil {
		log.Fatal(err)
	}
}
