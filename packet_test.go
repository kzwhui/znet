package znet

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestPacket_Marshal(t *testing.T) {
	packet := Packet{}
	packet.Data = []byte("1234")
	fmt.Printf("raw data=%v\n", packet.Data)
	var stream []byte
	var err error
	if stream, err = packet.Marshal(); err != nil {
		t.Errorf("err=%v", err)
	}

	hexString := hex.EncodeToString(stream)
	if hexString != "0000000831323334" {
		t.Errorf("not equal")
	}
}

func TestPacket_Unmarshal(t *testing.T) {
	hexData := "0000000831323334"
	streamData, err := hex.DecodeString(hexData)
	if err != nil {
		t.Errorf("fail to decode hex")
	}
	packet := Packet{}
	err = packet.Unmarshal([]byte(streamData))
	if err != nil {
		t.Errorf("fail to unmarshal")
	}
	fmt.Printf("packet data=%v", string(packet.Data))
	if string(packet.Data) != "1234" {
		t.Errorf("not equal")
	}

}

func TestPacket_Check(t *testing.T) {
}
