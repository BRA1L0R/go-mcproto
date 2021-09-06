package mcproto

import (
	"errors"
	"fmt"
	"net"

	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
)

// LoginDisconnectError rapresents the disconnection of the client during the login state,
// during the initialization process
type LoginDisconnectError struct {
	// Reason is the json encoded Reason for which the client has been
	// disconnected from the server
	Reason string
}

func (disconnectError *LoginDisconnectError) Error() string {
	return fmt.Sprintf("mcproto: server disconnected the client during initialization with message: %s", disconnectError.Reason)
}

// Connect initializes the connection to the server.
//
// host must have the format of "host:port" as a port has to be specified in order
// to open a connection. 25565 is not taken for granted
func (mc *Client) Connect(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	mc.connection = conn.(*net.TCPConn)

	return nil
}

// Initializes the connection to the server by sending
// the handshake packet and the login packet
//
// host is the server fqdn or ip address of the server, port is the uint16 port where the server is listening on
//
// username is the in-game username the client will send to the server during handshaking. Might differ from the actual
// in-game username as the server sends a confirmation of it after the login state.
func (mc *Client) Initialize(host string, port uint16, protocolVersion int32, username string) error {
	if err := mc.Connect(fmt.Sprintf("%s:%v", host, port)); err != nil {
		return err
	}

	// Create handshake packet with the latest protocol version
	// and state information
	// NOTE: serveraddress and serverport are not used by the
	// notchian server, but it's best practice to fill them
	hp := models.HandshakePacket{
		MinecraftPacket: packets.MinecraftPacket{PacketID: 0x00},

		ProtocolVersion: protocolVersion,
		ServerAddress:   host,
		ServerPort:      port,
		NextState:       2,
	}

	err := mc.WritePacket(&hp)
	if err != nil {
		return err
	}

	loginPacket := models.LoginStartPacket{
		MinecraftPacket: packets.MinecraftPacket{PacketID: 0x00},

		Name: username,
	}

	err = mc.WritePacket(&loginPacket)
	if err != nil {
		return err
	}

	for {
		p, err := mc.ReceivePacket()
		if err != nil {
			return err
		}

		switch p.PacketID {
		case 0x00: // disconnected
			disconnectPacket := new(models.DisconnectPacket)
			if err := p.DeserializeData(disconnectPacket); err != nil {
				return err
			}

			return &LoginDisconnectError{Reason: disconnectPacket.Reason}
		case 0x01: // encryption request
			return errors.New("mcproto: received an encryption request which means the server is in online mode. Online mode not currently supported")
		case 0x03: // set compression
			setCompression := new(models.SetCompressionPacket)
			err := p.DeserializeData(setCompression)
			if err != nil {
				return err
			}

			if setCompression.Treshold < 0 {
				return errors.New("mcproto: server sent a set compression packet with a negative treshold")
			}

			mc.CompressionTreshold = setCompression.Treshold
		case 0x02: // login success
			loginSuccess := new(models.LoginSuccessPacket)
			err := p.DeserializeData(loginSuccess)

			return err
		}
	}
}
