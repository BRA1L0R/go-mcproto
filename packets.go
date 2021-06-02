package mcprot

import (
	"bytes"
	"compress/zlib"
	"errors"
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

	// _, err = mc.connection.Read(packetContent)
	_, err = io.ReadFull(mc.connection, packetContent)
	if err != nil {
		return nil, err
	}

	packet := packets.NewUncompressedPacket(0)

	packet.Length = packetLength
	packet.PacketID = packetId
	packet.Data.Write(packetContent)

	return packet, nil
}

func (mc *McProt) ReceivePacket() (*packets.CompressedPacket, error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return nil, err
	}

	dataLength, dLenLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		// drain rest of the package to avoid problems
		return nil, err
	}

	remainingData := make([]byte, packetLength-dLenLen)
	read, err := io.ReadFull(mc.connection, remainingData)
	if err != nil {
		return nil, err
	}

	remainingDataBuffer := bytes.NewBuffer(remainingData)

	newPacket := packets.NewCompressedPacket(0)
	newPacket.Length = packetLength
	newPacket.DataLength = dataLength

	if dataLength != 0 {
		if int32(read) != (packetLength - dLenLen) {
			return nil, errors.New("bytes read from buffer and bytes that needed to be fulfilled mismatch")
		}

		reader, err := zlib.NewReader(remainingDataBuffer)
		defer reader.Close()
		if err != nil {
			return nil, err
		}

		uncompressedData := make([]byte, dataLength)
		_, err = reader.Read(uncompressedData)
		if err != nil && err != io.EOF {
			return nil, err
		}

		newPacket.Data = bytes.NewBuffer(uncompressedData)

		packetId, _, err := varint.DecodeReaderVarInt(newPacket.Data)
		if err != nil {
			return nil, err
		}

		newPacket.PacketID = packetId

		return newPacket, err
	} else {
		packetId, _, err := varint.DecodeReaderVarInt(remainingDataBuffer)
		if err != nil {
			return nil, err
		}

		newPacket.PacketID = packetId
		newPacket.Data = remainingDataBuffer

		return newPacket, nil
	}
}
