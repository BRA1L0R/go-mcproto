package mcproto

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"

	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/varint"
)

var (
	maxPacketLength    = 2097151 // https://wiki.vg/Protocol#Packet_format
	ErrMaxPacketLength = errors.New("mcproto: counterpart sent a packet which was too big")
)

func (mc *Client) ReceivePacket() (packet packets.MinecraftPacket, err error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return
	}

	if packetLength > int32(maxPacketLength) {
		return packet, ErrMaxPacketLength
	}

	// as wiki.vg, negative or zero values for compression will disable compression,
	// meaning the format will remain uncompressed
	if mc.CompressionTreshold <= 0 {
		return mc.receiveUncompressedPacket(packetLength)
	} else {
		return mc.receiveCompressedPacket(packetLength)
	}
}

func (mc *Client) receiveUncompressedPacket(packetLength int32) (packet packets.MinecraftPacket, err error) {
	packetId, packetIdLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		return packet, err
	}

	packet.PacketID = packetId

	if packetLength > int32(packetIdLen) {
		packetContent := make([]byte, packetLength-int32(packetIdLen))
		_, err = io.ReadFull(mc.connection, packetContent)
		if err != nil {
			return packet, err
		}

		packet.Data = bytes.NewBuffer(packetContent)
	}

	return
}

func (mc *Client) receiveCompressedPacket(packetLength int32) (packet packets.MinecraftPacket, err error) {
	dataLength, dLenLen, err := varint.DecodeReaderVarInt(mc.connection)
	if err != nil {
		// drain rest of the package to avoid problems
		return packet, err
	}

	remainingData := make([]byte, packetLength-int32(dLenLen))
	read, err := io.ReadFull(mc.connection, remainingData)
	if err != nil {
		return packet, err
	}

	remainingDataBuffer := bytes.NewBuffer(remainingData)

	if dataLength != 0 {
		// While the previous check in ReceivePacket() referred to the uncompressed length
		// of the whole packet, this check refers to the length of the uncompressed data
		if dataLength > int32(maxPacketLength) {
			return packet, ErrMaxPacketLength
		}

		if int32(read) != (packetLength - int32(dLenLen)) {
			return packet, errors.New("mcproto: bytes read from buffer and bytes that needed to be fulfilled mismatch")
		}

		reader, err := zlib.NewReader(remainingDataBuffer)
		if err != nil {
			return packet, err
		}

		uncompressedData := make([]byte, dataLength)
		_, err = io.ReadFull(reader, uncompressedData)
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
