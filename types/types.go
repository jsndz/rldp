package types

const (
	DATA = 1
	ACK  = 2
)

type Frame struct {
	Seq     uint32
	Ack     uint32
	Type    uint8
	Payload []byte
}
