package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
	"github.com/BRA1L0R/go-mcproto/varint"
)

var (
	host = flag.String("host", "127.0.0.1", "Server host")
	port = flag.Uint("port", 25565, "Server port")
)

func main() {
	flag.Parse()

	client := mcproto.McProto{
		Host:            *host,
		Port:            uint16(*port),
		ProtocolVersion: 754, // 1.16.5
		Name:            "GolangTest",
	}

	err := client.Initialize()
	if err != nil {
		panic(err)
	}

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
			receivedKeepalive := models.KeepAlivePacket{CompressedPacket: packet}

			err := receivedKeepalive.DeserializeData(&receivedKeepalive)
			if err != nil {
				panic(err)
			}

			serverBoundKeepalive := models.KeepAlivePacket{
				CompressedPacket: packets.NewCompressedPacket(0x10),
				KeepAliveID:      receivedKeepalive.KeepAliveID,
			}

			err = client.WritePacket(&serverBoundKeepalive)
			if err != nil {
				panic(err)
			}

			fmt.Println("KeepAlive sent")
		}
	}
}
