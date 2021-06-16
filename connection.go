package mcproto

import "net"

// SerializablePacket defines the standard methods that a struct should have
// in order to be serializable by the library
//
// You can actually create your own methods as long as they respect this standard
type SerializablePacket interface {
	// SerializeData takes an interface pointer as input and serializes all the fields in the
	// data buffer. It can and will return an error in case of invalid data
	SerializeData(inter interface{}) error

	// Serialize takes compressionTreshold as an input, which can be negative in case of a disabled
	// compression, and serializes it into a byte slice
	Serialize(compressionTreshold int32) ([]byte, error)
}

// WritePacket calls SerializeData and then calls WriteRawPacket
func (mc *Client) WritePacket(packet SerializablePacket) error {
	if err := packet.SerializeData(packet); err != nil {
		return err
	}

	return mc.WriteRawPacket(packet)
}

// WriteRawPacket takes a packet with PacketID and Data already filled out, and serializes
// it to a byte slice, which is subsequently written to the underlying connection
func (mc *Client) WriteRawPacket(packet SerializablePacket) error {
	serialized, err := packet.Serialize(mc.CompressionTreshold)
	if err != nil {
		return err
	}

	_, err = mc.connection.Write(serialized)
	return err
}

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
