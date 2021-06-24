package mcproto

import (
	"io"

	"github.com/BRA1L0R/go-mcproto/packets"
)

// SerializablePacket defines the standard methods that a struct should have
// in order to be serializable by the library
//
// You can actually create your own methods as long as they respect this standard
type SerializablePacket interface {
	// SerializeData takes an interface pointer as input and serializes all the fields in the
	// data buffer. It can and will return an error in case of invalid data
	SerializeData(inter interface{}) error

	SerializeCompressed(writer io.Writer, compressionTreshold int) error
	SerializeUncompressed(writer io.Writer) error
}

// WritePacket calls SerializeData and then calls WriteRawPacket
func (mc *Client) WritePacket(packet SerializablePacket) error {
	if err := packet.SerializeData(packet); err != nil {
		return err
	}

	if mc.IsCompressionEnabled() {
		return packet.SerializeCompressed(mc.connection, int(mc.CompressionTreshold))
	} else {
		return packet.SerializeUncompressed(mc.connection)
	}
}

// WriteRawPacket takes a rawpacket as input and serializes it in the connection
func (mc *Client) WriteRawPacket(rawPacket *packets.MinecraftRawPacket) error {
	if mc.IsCompressionEnabled() {
		return rawPacket.WriteCompressed(mc.connection)
	} else {
		return rawPacket.WriteUncompressed(mc.connection)
	}
}

// ReceiveRawPacket reads a raw packet from the connection but doesn't deserialize
// neither uncompress it
func (mc *Client) ReceiveRawPacket() (*packets.MinecraftRawPacket, error) {
	if mc.IsCompressionEnabled() {
		return packets.FromCompressedReader(mc.connection)
	} else {
		return packets.FromUncompressedReader(mc.connection)
	}
}

// ReceivePacket receives and deserializes a packet from the connection, uncompressing it
// if necessary
func (mc *Client) ReceivePacket() (*packets.MinecraftPacket, error) {
	rawPacket, err := mc.ReceiveRawPacket()
	if err != nil {
		return nil, err
	}

	return packets.FromRawPacket(rawPacket)
}
