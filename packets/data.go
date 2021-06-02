package packets

import (
	"bytes"

	"github.com/BRA1L0R/go-mcprot/packets/serialization"
)

type standardPacket struct {
	Data   *bytes.Buffer
	Length int32

	PacketID int32
}

func (p *standardPacket) SerializeData(inter interface{}) error {
	return serialization.SerializeFields(inter, p.Data)
}

func (p *standardPacket) DeserializeData(inter interface{}) error {
	return serialization.DeserializeFields(inter, p.Data)
}

func (p *standardPacket) InitializePacket() {
	p.Data = bytes.NewBuffer([]byte{})
}
