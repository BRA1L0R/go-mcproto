package mcprot

import (
	"fmt"
	"net"

	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/packets/models"
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
	hp := models.HandshakePacket{
		UncompressedPacket: packets.NewUncompressedPacket(0x00),

		ProtocolVersion: mc.ProtocolVersion,
		ServerAddress:   mc.Host,
		ServerPort:      mc.Port,
		NextState:       2,
	}

	err = mc.WritePacket(&hp)
	if err != nil {
		return err
	}

	loginPacket := models.LoginStartPacket{
		UncompressedPacket: packets.NewUncompressedPacket(0x00),

		Name: mc.Name,
	}

	err = mc.WritePacket(&loginPacket)
	if err != nil {
		return err
	}

	p, err := mc.ReceiveUncompressedPacket()
	if err != nil {
		return err
	}

	setCompPacket := models.SetCompressionPacket{UncompressedPacket: p}
	if setCompPacket.PacketID != 0x03 {
		panic(setCompPacket.Data.String())
	}

	err = setCompPacket.DeserializeData(&setCompPacket)
	if err != nil {
		return err
	}

	mc.compressionTreshold = setCompPacket.Treshold

	pack, err := mc.ReceivePacket()
	if err != nil {
		return err
	}

	loginSuccessPacket := models.LoginSuccessPacket{CompressedPacket: pack}
	if loginSuccessPacket.PacketID != 0x02 {
		panic(loginSuccessPacket.Data.String())
	}

	err = loginSuccessPacket.DeserializeData(&loginSuccessPacket)
	if err != nil {
		return err
	}

	return nil
}
