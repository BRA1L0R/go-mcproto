package models

import "github.com/BRA1L0R/go-mcproto/packets"

type LoginSuccessPacket struct {
	*packets.CompressedPacket

	UUID     string `type:"ignore" len:"16"`
	Username string `type:"string"`
}
