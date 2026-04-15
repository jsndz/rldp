package protocol

import (
	"bytes"
	"encoding/binary"

	"github.com/jsndz/rldp/types"
)

func Encoding(frame types.Frame) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, uint32(frame.Seq))
	binary.Write(&buf, binary.BigEndian, uint32(frame.Ack))
	buf.WriteByte(byte(frame.Type))
	buf.Write(frame.Payload)
	return buf.Bytes()
}

func Decoding(frame []byte) (uint32, uint32, uint8, string) {
	var seq uint32
	var ack uint32
	var frameType uint8

	buf := bytes.NewReader(frame)
	binary.Read(buf, binary.BigEndian, &seq)
	binary.Read(buf, binary.BigEndian, &ack)
	frameTypeByte, _ := buf.ReadByte()
	frameType = uint8(frameTypeByte)

	payload := make([]byte, buf.Len())
	buf.Read(payload)

	return seq, ack, frameType, string(payload)
}
