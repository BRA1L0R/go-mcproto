package mcproto

import (
	"net"
)

type McProto struct {
	Host            string
	Port            uint16
	Name            string
	ProtocolVersion int32

	compressionTreshold int32

	connection net.Conn
}

func (mc *McProto) GetCompressionTreshold() int32 {
	return mc.compressionTreshold
}
