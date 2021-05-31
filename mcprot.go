package mcprot

import (
	"net"
)

type McProt struct {
	Host string
	Port uint16
	Name string

	CompressionTreshold int

	Connection net.Conn
}
