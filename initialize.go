package mcprot

import (
	"fmt"
	"net"

	"github.com/BRA1L0R/go-mcprot/packets"
)

func (mc *McProt) Initialize() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", mc.Host, mc.Port))
	mc.Connection = conn
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

		ProtocolVersion: 754,
		ServerAddress:   mc.Host,
		ServerPort:      mc.Port,
		NextState:       2,
	}

	if err := hp.SerializeData(hp); err != nil {
		return err
	}

	conn.Write(hp.Serialize())

	loginPacket := packets.LoginStartPacket{
		UncompressedPacket: &packets.UncompressedPacket{
			PacketID: 0x00,
		},

		Name: mc.Name,
	}

	if err := loginPacket.SerializeData(loginPacket); err != nil {
		return err
	}

	conn.Write(loginPacket.Serialize())

	p, err := mc.ReceiveUncompressedPacket()
	if err != nil {
		return err
	}

	setCompPacket := packets.SetCompressionPacket{UncompressedPacket: p}
	if setCompPacket.PacketID != 0x03 {
		panic(setCompPacket.Data.String())
	}

	setCompPacket.DeserializeData(&setCompPacket)
	mc.CompressionTreshold = setCompPacket.Treshold

	pack, err := mc.ReceivePacket()
	if err != nil {
		return err
	}

	loginSuccessPacket := packets.LoginSuccessPacket{StandardPacket: pack}
	if loginSuccessPacket.PacketID != 0x02 {
		panic(loginSuccessPacket.Data.String())
	}

	loginSuccessPacket.DeserializeData(&loginSuccessPacket)

	fmt.Println("Logged in with username: " + loginSuccessPacket.Username)
	return nil
}
