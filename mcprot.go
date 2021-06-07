package mcproto

import (
	"net"
)

type Client struct {
	// IP or FQDN the server will try to enstablish a connection to
	Host            string

	// Port which the server listens on
	Port            uint16

	// In-game username of the connecting client 
	Name            string

	// Protocol version implemented by the client.
	// All the protocol versions are available here: https://wiki.vg/Protocol_version_numbers
	ProtocolVersion int32

	// determines whenever a packet should be compressed or not
	// if negative, the library will serialize packets using the Uncompressed Format (wiki.vg)
	//
	// automatically set by the Client.Initialize method
	CompressionTreshold int32

	connection net.Conn
}