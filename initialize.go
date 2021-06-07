package mcproto

import (
	"errors"
	"fmt"
	"net"

	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
)

func (mc *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", mc.Host, mc.Port))
	mc.connection = conn

	return err
}

// Initializes the connection to the server by sending
// the handshake packet and the login packet
//
// Server Host, Port and Username are defined in the McProto object
func (mc *Client) Initialize() error {
	if err := mc.Connect(); err != nil {
		return err
	}

	// Create handshake packet with the latest protocol version
	// and state information
	// NOTE: serveraddress and serverport are not used by the
	// notchian server, but it's best practice to fill them
	hp := models.HandshakePacket{
		MinecraftPacket: packets.MinecraftPacket{PacketID: 0x00},

		ProtocolVersion: mc.ProtocolVersion,
		ServerAddress:   mc.Host,
		ServerPort:      mc.Port,
		NextState:       2,
	}

	err := mc.WritePacket(&hp)
	if err != nil {
		return err
	}

	loginPacket := models.LoginStartPacket{
		MinecraftPacket: packets.MinecraftPacket{PacketID: 0x00},

		Name: mc.Name,
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
			disconnectPacket := models.DisconnectPacket{MinecraftPacket: p}
			if err := disconnectPacket.DeserializeData(&disconnectPacket); err != nil {
				return err
			}

			return fmt.Errorf("mcproto: the server disconnected the client, reason: %s", disconnectPacket.Reason)
		case 0x01: // encryption request
			return errors.New("mcproto: received an encryption request which means the server is in online mode. Online mode not currently supported")
		case 0x03: // set compression
			setCompression := models.SetCompressionPacket{MinecraftPacket: p}
			err := setCompression.DeserializeData(&setCompression)
			if err != nil {
				return err
			}

			if setCompression.Treshold < 0 {
				return errors.New("mcproto: server sent a set compression packet with a negative treshold")
			}

			mc.CompressionTreshold = setCompression.Treshold
		case 0x02: // login success
			loginSuccess := models.LoginStartPacket{MinecraftPacket: p}
			err := loginSuccess.DeserializeData(&loginSuccess)

			return err
		}
	}
}
