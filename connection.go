package mcproto

import "net"

// FromListener accepts the connection from a net.Listener and uses it as the
// underlying connection for packet comunication between the server and the client
func FromListener(ln *net.TCPListener) (client *Client, err error) {
	client = new(Client)
	client.connection, err = ln.AcceptTCP()

	return
}

// CloseConnection closes the underlying connection making further comunication impossible
func (mc *Client) CloseConnection() error {
	return mc.connection.Close()
}
