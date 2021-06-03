package models

import "github.com/BRA1L0R/go-mcproto/packets"

type HandshakePacket struct {
	// *packets.UncompressedPacket
	packets.MinecraftPacket

	ProtocolVersion int32  `type:"varint"`
	ServerAddress   string `type:"string"`
	ServerPort      uint16 `type:"inherit"`
	NextState       int32  `type:"varint"`
}
