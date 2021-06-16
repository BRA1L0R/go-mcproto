package mcproto

import (
	"errors"
	"fmt"
	"net"

	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
)

func (mc *Client) Connect(host string, port uint16) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))

	mc.connection = conn.(*net.TCPConn)

	return err
}

// Initializes the connection to the server by sending
// the handshake packet and the login packet
//
// host is the server fqdn or ip address of the server, port is the uint16 port where the server is listening on
//
// username is the in-game username the client will send to the server during handshaking. Might differ from the actual
// in-game username as the server sends a confirmation of it after the login state.
func (mc *Client) Initialize(host string, port uint16, protocolVersion int32, username string) error {
	if err := mc.Connect(host, port); err != nil {
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

			return fmt.Errorf("mcproto: the server disconnected the client, reason: %s", disconnectPacket.Reason)
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
