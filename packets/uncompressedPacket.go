package packets

import (
	"bytes"

	"github.com/BRA1L0R/go-mcproto/varint"
)

type UncompressedPacket struct {
	*standardPacket
}

// Serialize() on UncompressedPackets never returns an error
// but it's in the return params to be compliant to the PacketOp interface
//
// note: compressionTreshold is not used in uncompressed packets
func (p *UncompressedPacket) Serialize(compressionTreshold int32) ([]byte, error) {
	encodedPacketId, pIdLength := varint.EncodeVarInt(p.PacketID)

	p.Length = pIdLength + int32(p.Data.Len())
	encodedLength, _ := varint.EncodeVarInt(p.Length)

	dataBuffer := new(bytes.Buffer)
	dataBuffer.Write(encodedLength)
	dataBuffer.Write(encodedPacketId)

	_, err := dataBuffer.ReadFrom(p.Data)
	if err != nil {
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

func NewUncompressedPacket(packetId int32) *UncompressedPacket {
	uncompressedPacket := new(UncompressedPacket)
	uncompressedPacket.standardPacket = new(standardPacket)
	uncompressedPacket.InitializePacket()
	uncompressedPacket.PacketID = packetId

	return uncompressedPacket
}
