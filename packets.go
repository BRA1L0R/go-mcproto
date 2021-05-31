package mcprot

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/varint"
)

func (mc *McProt) ReceiveUncompressedPacket() (*packets.UncompressedPacket, error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return nil, err
	}

	packetId, packetIdLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return nil, err
	}

	packetContent := make([]byte, packetLength-packetIdLen)
	mc.connection.Read(packetContent)

	packet := new(packets.UncompressedPacket)

	packet.Length = packetLength
	packet.PacketID = packetId
	packet.Data.Write(packetContent)

	return packet, nil
}

func (mc *McProt) ReceivePacket() (*packets.StandardPacket, error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return nil, err
	}

	dataLength, dLenLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		// drain rest of the package to avoid problems
		mc.drain(packetLength - dLenLen)
		return nil, err
	}

	if dataLength != 0 {
		// drain the remaining packet
		compressedData := make([]byte, packetLength-dLenLen)
		mc.connection.Read(compressedData)

		compressedBuffer := bytes.NewBuffer(compressedData)

		reader, err := zlib.NewReader(compressedBuffer)
		if err != nil {
			return nil, err
		}

		uncompressedBuffer := new(bytes.Buffer)
		io.Copy(uncompressedBuffer, reader)

		packetId, pIdLen, err := varint.DecodeReaderVarInt(uncompressedBuffer)
		if err != nil {
			mc.drain(dataLength - pIdLen)
			return nil, err
		}

		newPacket := new(packets.StandardPacket)
		newPacket.Length = packetLength
		newPacket.DataLength = dataLength
		newPacket.PacketID = packetId
		newPacket.Data.ReadFrom(uncompressedBuffer)

		return newPacket, nil
	} else {
		packetId, pIdLen, err := varint.DecodeReaderVarInt(mc.connection)
		if err != nil {
			mc.drain(packetLength - (dLenLen + pIdLen))
			return nil, err
		}

		newPacket := new(packets.StandardPacket)
		newPacket.Length = packetLength
		newPacket.DataLength = dataLength
		newPacket.PacketID = packetId

		remainingDataLen := packetLength - dLenLen - pIdLen

		data := make([]byte, remainingDataLen)
		mc.connection.Read(data)

		newPacket.Data.Write(data)

		return newPacket, nil
	}
}

func (mc *McProt) drain(size int) {
	drain := make([]byte, size)
	mc.connection.Read(drain)
}
