package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PacketInterface interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	// < 0 -- 包长度不够
	// >= 0 -- 包长度超出多少
	Check([]byte) int64
}

// 4字节表示总长度，后面接着正文
type Packet struct {
	Data []byte
}

func (p *Packet) Marshal() (stream []byte, err error) {
	lengthBytes := make([]byte, 4)
	length := uint32(len(p.Data) + 4)
	binary.BigEndian.PutUint32(lengthBytes[0:], length)

	var buffer bytes.Buffer
	buffer.Write(lengthBytes)
	buffer.Write(p.Data)
	stream = buffer.Bytes()

	return
}

func (p *Packet) Unmarshal(stream []byte) (err error) {
	ret := p.Check(stream)
	if ret < 0 {
		err = fmt.Errorf("packet is not enough")
		return
	}

	if ret > 0 {
		err = fmt.Errorf("packet has extra data")
		return
	}

	p.Data = stream[4:]
	return
}

// < 0 -- 包长度不够
// >= 0 -- 包长度超出多少
func (p *Packet) Check(stream []byte) int64 {
	streamLen := int64(len(stream))
	if len(stream) < 5 {
		return -1
	}

	total := int64(binary.BigEndian.Uint32(stream[0:4]))
	if total > streamLen {
		return -1
	}

	return streamLen - total
}
