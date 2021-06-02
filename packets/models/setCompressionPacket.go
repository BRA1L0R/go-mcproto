package models

import "github.com/BRA1L0R/go-mcprot/packets"

type SetCompressionPacket struct {
	*packets.UncompressedPacket

	Treshold int32 `type:"varint"`
}