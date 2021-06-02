package mcprot

import (
	"net"
)

type McProt struct {
	Host            string
	Port            uint16
	Name            string
	ProtocolVersion int32

	compressionTreshold int32

	connection net.Conn
}

func (mc *McProt) GetCompressionTreshold() int32 {
	return mc.compressionTreshold
}
