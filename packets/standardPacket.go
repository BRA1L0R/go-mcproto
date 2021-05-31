package packets

import (
	"bytes"
	"errors"

	"github.com/BRA1L0R/go-mcprot/varint"
)

type StandardPacket struct {
	Length int

	DataLength int

	// Rest is either compressed or uncompressed depending on DataLength
	PacketID int
	Data     bytes.Buffer
}

func (p *StandardPacket) SerializeData(inter interface{}) error {
	return encodeFields(inter, &p.Data)
}

func (p *StandardPacket) DeserializeData(inter interface{}) error {
	return decodeFields(inter, &p.Data)
}

func (p *StandardPacket) Serialize(compressionTreshold int) ([]byte, error) {
	dataSize := p.Data.Len()
	if dataSize >= compressionTreshold {
		return nil, errors.New("serialization with compression not implemented")
	}

	packetBuffer := new(bytes.Buffer)

	packetID, packetIDSize := varint.EncodeVarInt(p.PacketID)
	dataLength, dataLengthSize := varint.EncodeVarInt(0) // data length must be 0 in uncompressed packets
	packetLength, _ := varint.EncodeVarInt(packetIDSize + dataLengthSize + dataSize)

	packetBuffer.Write(packetLength)
	packetBuffer.Write(dataLength)
	packetBuffer.Write(packetID)
	packetBuffer.ReadFrom(&p.Data)

	return packetBuffer.Bytes(), nil
}
