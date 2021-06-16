package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/packets/models"
	"github.com/BRA1L0R/go-mcproto/varint"
)

var (
	host     = flag.String("host", "127.0.0.1", "Server host")
	port     = flag.Uint("port", 25565, "Server port")
	username = flag.String("username", "GolangKeepalive", "In-game username")
)

func main() {
	flag.Parse()

	client := mcproto.Client{}
	err := client.Initialize(*host, uint16(*port), 754, *username)
	if err != nil {
		panic(err)
	}

	fmt.Println("Joined the server as:", *username)

	for {
		packet, err := client.ReceivePacket()
		if err != nil {
			if err == varint.ErrVarIntTooBig {
				fmt.Println("Received a varint which was too big")
			}
			continue
		}

		if packet.PacketID == 0x1F {
			receivedKeepalive := new(models.KeepAlivePacket)
			err := packet.DeserializeData(receivedKeepalive)
			if err != nil {
				panic(err)
			}

			serverBoundKeepalive := new(models.KeepAlivePacket)
			serverBoundKeepalive.KeepAliveID = receivedKeepalive.KeepAliveID
			serverBoundKeepalive.PacketID = 0x10

			err = client.WritePacket(serverBoundKeepalive)
			if err != nil {
				panic(err)
			}

			fmt.Println("KeepAlive sent")
		}
	}
}
