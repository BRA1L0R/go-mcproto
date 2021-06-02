package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcprot"
	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/varint"
)

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
		Name:            "GolangTest",
	}

	client.Initialize()

	fmt.Println("Joined the server as:", client.Name)

	for {
		packet, err := client.ReceivePacket()
		if err != nil {
			if err == varint.ErrVarIntTooBig {
				fmt.Println("Received a varint which was too big")
			}
			continue
		}

		if packet.PacketID == 0x1F {
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
			fmt.Println("KeepAlive sent")
		}
	}
}
