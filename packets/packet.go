package packets

import (
	"bytes"
	"io"
)

// Basic struct containing the PacketID and many methods
// for serialization.
//
// Compliant to the PacketOP interface
type MinecraftPacket struct {
	// PacketID is a VarInt identifier of the packet
	// refer to wiki.vg for all the packet ids
	PacketID int32

	// Contains inflated data. Never contains compressed zlib data.
	// Check if data has been inflated and is present in this buffer
	// using IsCompressed()
	Data *bytes.Buffer

	// Indicates whenever or not the Data buffer is a compressed
	// zlib buffer or it's already in its uncompressed form
	compressed bool
	// zlib reader
	compressedData io.ReadCloser
}

func (p *MinecraftPacket) IsCompressed() bool {
	return p.compressed
}

// When compressionTreshold is a value below zero
// the packet gets serialized with the "Without Compression"
// format. This happens when a server never sent a set compression
// packet
func (p *MinecraftPacket) Serialize(compressionTreshold int32) ([]byte, error) {
	if compressionTreshold <= 0 {
		return p.serializeUncompressedPacket()
	} else {
		return p.serializeCompressedPacket(compressionTreshold)
	}
}

// func NewPacket() *MinecraftPacket {
// 	return new(MinecraftPacket)
// }

func NewPacketWithCompressedData(PacketId int32, compressedData io.ReadCloser) *MinecraftPacket {
	packet := new(MinecraftPacket)
	packet.PacketID = PacketId

	packet.compressedData = compressedData
	packet.compressed = true

	return packet
}

func NewPacketWithData(PacketId int32, data *bytes.Buffer) *MinecraftPacket {
	packet := new(MinecraftPacket)
	packet.PacketID = PacketId

	packet.Data = data
	packet.compressed = false

	return packet
}
