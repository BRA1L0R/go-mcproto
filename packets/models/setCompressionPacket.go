package models

import "github.com/BRA1L0R/go-mcproto/packets"

type SetCompressionPacket struct {
	// *packets.UncompressedPacket
	packets.MinecraftPacket

	Treshold int32 `mc:"varint"`
}
