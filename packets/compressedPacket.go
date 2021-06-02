package packets

import (
	"bytes"
	"errors"

	"github.com/BRA1L0R/go-mcprot/varint"
)

type CompressedPacket struct {
	*standardPacket
	DataLength int32
}

func (p *CompressedPacket) Serialize(compressionTreshold int32) ([]byte, error) {
	dataSize := int32(p.Data.Len())
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

	_, err := packetBuffer.ReadFrom(p.Data)
	if err != nil {
		return nil, err
	}

	return packetBuffer.Bytes(), nil
}

func NewCompressedPacket(packetId int32) *CompressedPacket {
	compressedPacket := new(CompressedPacket)
	compressedPacket.standardPacket = new(standardPacket)
	compressedPacket.InitializePacket()
	compressedPacket.PacketID = packetId

	return compressedPacket
}
