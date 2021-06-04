package models

import "github.com/BRA1L0R/go-mcproto/packets"

type KeepAlivePacket struct {
	// *packets.CompressedPacket
	packets.MinecraftPacket

	KeepAliveID int64 `mc:"inherit"`
}
