package mcprot

import (
	"net"
)

type McProt struct {
	Host            string
	Port            uint16
	Name            string
	ProtocolVersion int

	compressionTreshold int

	connection net.Conn
}

func (mc *McProt) GetCompressionTreshold() int {
	return mc.compressionTreshold
}
