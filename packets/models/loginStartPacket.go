package models

import "github.com/BRA1L0R/go-mcproto/packets"

type LoginStartPacket struct {
	*packets.UncompressedPacket

	Name string `type:"string"`
}
