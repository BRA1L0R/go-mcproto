package models

import "github.com/BRA1L0R/go-mcproto/packets"

type LoginSuccessPacket struct {
	// *packets.CompressedPacket
	packets.MinecraftPacket

	UUID     []byte `mc:"bytes" len:"16"`
	Username string `mc:"string"`
}
