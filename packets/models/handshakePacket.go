package models

import "github.com/BRA1L0R/go-mcproto/packets"

type HandshakePacket struct {
	// *packets.UncompressedPacket
	packets.MinecraftPacket

	ProtocolVersion int32  `mc:"varint"`
	ServerAddress   string `mc:"string"`
	ServerPort      uint16 `mc:"inherit"`
	NextState       int32  `mc:"varint"`
}
