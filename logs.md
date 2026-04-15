Lets start with simple idea:

In terminal type send hello to send message
the have a goroutine running in background to listen to messages

App → TCP → IP → Ethernet → NIC → wire

Your case:

App → Ethernet (your code) → NIC → wire

you connect to raw socket.

socket is the thing that connect software to network hardware.
direct access to layer 2
you build the frame
send the frame
os gives it to NIC
NIC puts in the network
Switch/router sees destination MAC
NIC of other VM receives the frame

MAC is address in layer 2

With raw socket:

you bypass most of the stack
you directly send/receive Ethernet frames

NIC (Network Interface Card) = hardware that connects your machine to a network.

nic converts it to electric or radio signals

in simple words send works like:

get your inferface which gives you mac address
get dest mac address
give a unique bit for your etherType
so that you can tell the reciever we are using this protocol
get you payload

convert all this into a frame

form the address and send it to the Unix socket which connects to hardware
using sendto

Sender code snippet:

```go

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
```

REciever:::

A network interface is a system’s network endpoint (physical or virtual).

Examples:

eth0 → Ethernet card
wlan0 → Wi-Fi
lo → loopback

we create a socket and we need to bind it to one so that it will be only fixed to one
like we are binding to eth0 here

create socket which return a file descriptor
make an addr for eth0 in LinkLayer format

the bind the socket for that eth0

then when you recieve form that socket print it

Receiver code snippet:

```go


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

func Receive() {
	println("receiving data")
	iface, _ := net.InterfaceByName("wlan0")

	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(unix.ETH_P_ALL))
	
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
```



The testing with raw ethernet was hard and not working so switched to UDP
also added sq,ack, type and payload to the packet