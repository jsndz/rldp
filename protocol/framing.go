package protocol

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"

	"github.com/jsndz/rldp/types"
)

func Encoding(frame types.Frame) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, uint32(frame.Seq))
	binary.Write(&buf, binary.BigEndian, uint32(frame.Ack))
	buf.WriteByte(byte(frame.Type))
	buf.Write(frame.Payload)
	checksum := crc32.ChecksumIEEE(buf.Bytes())
	binary.Write(&buf, binary.BigEndian, checksum)
	return buf.Bytes()
}

func Decoding(frame []byte) (uint32, uint32, uint8, string, uint32) {
	var seq uint32
	var ack uint32
	var frameType uint8
	var checksum uint32
	buf := bytes.NewReader(frame)
	binary.Read(buf, binary.BigEndian, &seq)
	binary.Read(buf, binary.BigEndian, &ack)
	frameTypeByte, _ := buf.ReadByte()
	frameType = uint8(frameTypeByte)
	if buf.Len() < 4 {
		return 0, 0, 0, "", 0
	}
	remainingBytes := buf.Len()
	payload := make([]byte, remainingBytes-4)
	buf.Read(payload)
	binary.Read(buf, binary.BigEndian, &checksum)
	return seq, ack, frameType, string(payload), checksum
}

func VerifyChecksum(frame []byte) bool {
	if len(frame) < 9 {
		return false
	}
	data := frame[:len(frame)-4]
	checksumBytes := frame[len(frame)-4:]
	var checksum uint32
	buf := bytes.NewReader(checksumBytes)
	binary.Read(buf, binary.BigEndian, &checksum)
	calculatedChecksum := crc32.ChecksumIEEE(data)
	return checksum == calculatedChecksum
}
