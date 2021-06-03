package models

import "github.com/BRA1L0R/go-mcproto/packets"

type DisconnectPacket struct {
	packets.MinecraftPacket

	Reason string `type:"string"`
}
