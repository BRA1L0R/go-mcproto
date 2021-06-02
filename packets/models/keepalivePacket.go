package models

import "github.com/BRA1L0R/go-mcproto/packets"

type KeepAlivePacket struct {
	*packets.CompressedPacket

	KeepAliveID int64 `type:"inherit"`
}
