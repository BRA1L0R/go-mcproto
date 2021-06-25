package mcproto

import (
	"net"
	"sync"
)

type Client struct {
	// Determines whenever a packet should be compressed or not
	// if negative, the library will serialize packets using the Uncompressed Format (wiki.vg)
	//
	// automatically set by the Client.Initialize method, so if you are using this method
	// to open up the connection between the server and the client you must not modify it
	CompressionTreshold int32

	// read mutex to make the reading of packets thread safe
	readMu sync.Mutex
	// write mutex to make the writing of packets thread safe
	writeMu sync.Mutex

	// underlying tcp connection
	connection *net.TCPConn
}

func (client *Client) IsCompressionEnabled() bool {
	return client.CompressionTreshold > 0
}

func (client *Client) GetRemoteAddress() net.TCPAddr {
	return *client.connection.RemoteAddr().(*net.TCPAddr)
}
