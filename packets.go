package mcproto

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"

	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/varint"
)

func (mc *McProto) ReceivePacket() (packets.MinecraftPacket, error) {
	// as wiki.vg, negative or zero values for compression will disable compression,
	// meaning the format will remain uncompressed
	if mc.compressionTreshold <= 0 {
		return mc.receiveUncompressedPacket()
	} else {
		return mc.receiveCompressedPacket()
	}
}

func (mc *McProto) receiveUncompressedPacket() (packet packets.MinecraftPacket, err error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return packet, err
	}

	packetId, packetIdLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return packet, err
	}

	packetContent := make([]byte, packetLength-packetIdLen)
	_, err = io.ReadFull(mc.connection, packetContent)
	if err != nil {
		return packet, err
	}

	packet.PacketID = packetId
	packet.Data = bytes.NewBuffer(packetContent)

	return
}

func (mc *McProto) receiveCompressedPacket() (packet packets.MinecraftPacket, err error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return packet, err
	}

	dataLength, dLenLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		// drain rest of the package to avoid problems
		return packet, err
	}

	remainingData := make([]byte, packetLength-dLenLen)
	read, err := io.ReadFull(mc.connection, remainingData)
	if err != nil {
		return packet, err
	}

	remainingDataBuffer := bytes.NewBuffer(remainingData)

	if dataLength != 0 {
		if int32(read) != (packetLength - dLenLen) {
			return packet, errors.New("mcproto: bytes read from buffer and bytes that needed to be fulfilled mismatch")
		}

		reader, err := zlib.NewReader(remainingDataBuffer)
		if err != nil {
			return packet, err
		}

		uncompressedData := make([]byte, dataLength)
		_, err = reader.Read(uncompressedData)
		if err != nil && err != io.EOF {
			return packet, err
		}

		packet.Data = bytes.NewBuffer(uncompressedData)

		packetId, _, err := varint.DecodeReaderVarInt(packet.Data)
		if err != nil {
			return packet, err
		}

		packet.PacketID = packetId

		return packet, reader.Close()
	} else {
		packetId, _, err := varint.DecodeReaderVarInt(remainingDataBuffer)
		if err != nil {
			return packet, err
		}

		packet.PacketID = packetId
		packet.Data = remainingDataBuffer

		return packet, nil
	}
}
