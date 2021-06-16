package mcproto

import (
	"net"
)

type Client struct {
	// IP or FQDN the server will try to enstablish a connection to
	Host string

	// Port which the server listens on
	Port uint16

	// In-game username of the connecting client
	Name string

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

// FromListener accepts the connection from a net.Listener and uses it as the
// underlying connection for packet comunication between the server and the client
func FromListener(ln net.Listener) (client *Client, err error) {
	client = new(Client)

	client.connection, err = ln.Accept()
	client.Host = client.connection.RemoteAddr().String()
	return
}
