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

	client := mcprot.McProt{Host: *host, Port: uint16(*port), Name: "GolangTest"}
	client.Initialize()

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
				StandardPacket: &packets.StandardPacket{PacketID: 0x10},
				KeepAliveID:    receivedKeepalive.KeepAliveID,
			}

			err = serverBoundKeepalive.SerializeData(serverBoundKeepalive)
			if err != nil {
				panic(err)
			}

			keepAlive, err := serverBoundKeepalive.Serialize(client.CompressionTreshold)
			if err != nil {
				panic(err)
			}

			client.Connection.Write(keepAlive)
			fmt.Println("KeepAlive sent")
		}
	}
}
