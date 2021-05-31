package mcprot

import (
	"fmt"
	"net"

	"github.com/BRA1L0R/go-mcprot/packets"
)

func (mc *McProt) Initialize() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", mc.Host, mc.Port))
	mc.connection = conn
	if err != nil {
		return err
	}

	// Create handshake packet with the latest protocol version
	// and state information
	// NOTE: serveraddress and serverport are not used by the
	// notchian server, but it's best practice to fill them
	hp := packets.HandshakePacket{
		UncompressedPacket: &packets.UncompressedPacket{
			PacketID: 0x00,
		},

		ProtocolVersion: mc.ProtocolVersion,
		ServerAddress:   mc.Host,
		ServerPort:      mc.Port,
		NextState:       2,
	}

	mc.WriteUncompressedPacket(hp)

	loginPacket := packets.LoginStartPacket{
		UncompressedPacket: &packets.UncompressedPacket{
			PacketID: 0x00,
		},

		Name: mc.Name,
	}

	mc.WriteUncompressedPacket(loginPacket)

	p, err := mc.ReceiveUncompressedPacket()
	if err != nil {
		return err
	}

	setCompPacket := packets.SetCompressionPacket{UncompressedPacket: p}
	if setCompPacket.PacketID != 0x03 {
		panic(setCompPacket.Data.String())
	}

	setCompPacket.DeserializeData(&setCompPacket)
	mc.compressionTreshold = setCompPacket.Treshold

	pack, err := mc.ReceivePacket()
	if err != nil {
		return err
	}

	loginSuccessPacket := packets.LoginSuccessPacket{StandardPacket: pack}
	if loginSuccessPacket.PacketID != 0x02 {
		panic(loginSuccessPacket.Data.String())
	}

	loginSuccessPacket.DeserializeData(&loginSuccessPacket)
	return nil
}
