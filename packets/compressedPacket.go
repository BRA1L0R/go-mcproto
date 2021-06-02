package packets

import (
	"bytes"
	"compress/zlib"

	"github.com/BRA1L0R/go-mcprot/varint"
)

type CompressedPacket struct {
	*standardPacket
	DataLength int32
}

func (p *CompressedPacket) Serialize(compressionTreshold int32) ([]byte, error) {
	packetBuffer := new(bytes.Buffer)

	dataSize := int32(p.Data.Len())
	if dataSize >= compressionTreshold {
		err := p.serializeWithCompression(dataSize, packetBuffer)
		if err != nil {
			return nil, err
		}
	} else {
		err := p.serializeWithoutCompression(dataSize, packetBuffer)
		if err != nil {
			return nil, err
		}
	}

	return packetBuffer.Bytes(), nil
}

func (p *CompressedPacket) serializeWithoutCompression(dataSize int32, packetBuffer *bytes.Buffer) error {
	packetID, packetIDSize := varint.EncodeVarInt(p.PacketID)
	dataLength, dataLengthSize := varint.EncodeVarInt(0) // data length must be 0 in uncompressed packets
	packetLength, _ := varint.EncodeVarInt(packetIDSize + dataLengthSize + dataSize)

	packetBuffer.Write(packetLength)
	packetBuffer.Write(dataLength)
	packetBuffer.Write(packetID)

	_, err := packetBuffer.ReadFrom(p.Data)
	if err != nil {
		return err
	}

	return nil
}

func (p *CompressedPacket) serializeWithCompression(dataSize int32, packetBuffer *bytes.Buffer) error {
	packetID, packetIdSize := varint.EncodeVarInt(p.PacketID)
	dataLength, dataLengthSize := varint.EncodeVarInt(packetIdSize + dataSize)

	compressedBuf := new(bytes.Buffer)
	writer := zlib.NewWriter(compressedBuf)

	_, err := writer.Write(packetID)
	if err != nil {
		return err
	}

	_, err = writer.Write(p.Data.Bytes())
	if err != nil {
		return err
	}

	writer.Flush()
	writer.Close()

	packetLength, _ := varint.EncodeVarInt(dataLengthSize + int32(compressedBuf.Len()))

	packetBuffer.Write(packetLength)
	packetBuffer.Write(dataLength)

	_, err = packetBuffer.ReadFrom(compressedBuf)
	if err != nil {
		return err
	}

	return nil
}

func NewCompressedPacket(packetId int32) *CompressedPacket {
	compressedPacket := new(CompressedPacket)
	compressedPacket.standardPacket = new(standardPacket)
	compressedPacket.InitializePacket()
	compressedPacket.PacketID = packetId

	return compressedPacket
}
