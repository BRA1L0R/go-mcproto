package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcprot"
	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/varint"
)

type UpdateHealth struct {
	*packets.CompressedPacket

	Health         float32 `type:"inherit"`
	Food           int     `type:"varint"`
	FoodSaturation float32 `type:"inherit"`
}

var (
	host = flag.String("host", "127.0.0.1", "Server host")
	port = flag.Uint("port", 25565, "Server port")
)

func main() {
	flag.Parse()

	client := mcprot.McProt{
		Host:            *host,
		Port:            uint16(*port),
		ProtocolVersion: 754, // 1.16.5
		Name:            "GoHealth-Test",
	}

	client.Initialize()

	for {
		packet, err := client.ReceivePacket()
		if err != nil {
			if err == varint.ErrVarIntTooBig {
				fmt.Println("Received a varint which was too big")
			}
			continue
		}

		switch packet.PacketID { // Update Health
		case 0x49:
			healthPacket := UpdateHealth{StandardPacket: packet}
			healthPacket.DeserializeData(&healthPacket)

			fmt.Printf("Health Update: %vH %vF\n", healthPacket.Health, healthPacket.Food)
		case 0x1F:
			receivedKeepalive := packets.KeepAlivePacket{StandardPacket: packet}

			err := receivedKeepalive.DeserializeData(&receivedKeepalive)
			if err != nil {
				panic(err)
			}

			serverBoundKeepalive := packets.KeepAlivePacket{
				StandardPacket: &packets.CompressedPacket{PacketID: 0x10},
				KeepAliveID:    receivedKeepalive.KeepAliveID,
			}

			client.WritePacket(serverBoundKeepalive)
		}
	}
}
