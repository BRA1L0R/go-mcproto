package packets

import (
	"bytes"

	"github.com/BRA1L0R/go-mcprot/varint"
)

type UncompressedPacket struct {
	Length int

	PacketID int
	Data     bytes.Buffer
}

func (p *UncompressedPacket) Serialize() []byte {
	encodedPacketId, pIdLength := varint.EncodeVarInt(p.PacketID)

	p.Length = pIdLength + p.Data.Len()
	encodedLength, _ := varint.EncodeVarInt(p.Length)

	dataBuffer := new(bytes.Buffer)
	dataBuffer.Write(encodedLength)
	dataBuffer.Write(encodedPacketId)
	dataBuffer.ReadFrom(&p.Data)

	return dataBuffer.Bytes()
}

func (p *UncompressedPacket) SerializeData(inter interface{}) error {
	return encodeFields(inter, &p.Data)
}

func (p *UncompressedPacket) DeserializeData(inter interface{}) error {
	return decodeFields(inter, &p.Data)
}
