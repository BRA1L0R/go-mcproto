package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
	"github.com/BRA1L0R/go-mcproto/varint"
)

type UpdateHealth struct {
	packets.MinecraftPacket

	Health         float32 `mc:"inherit"`
	Food           int     `mc:"varint"`
	FoodSaturation float32 `mc:"inherit"`
}

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
			healthPacket := new(UpdateHealth)
			err := packet.DeserializeData(healthPacket)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Health Update: %vH %vF\n", healthPacket.Health, healthPacket.Food)
		case 0x1F:
			receivedKeepalive := new(models.KeepAlivePacket)
			err := packet.DeserializeData(receivedKeepalive)
			if err != nil {
				panic(err)
			}

			// serverBoundKeepalive := models.KeepAlivePacket{
			// 	CompressedPacket: packets.NewCompressedPacket(0x10),
			// 	KeepAliveID:      receivedKeepalive.KeepAliveID,
			// }
			serverBoundKeepalive := new(models.KeepAlivePacket)
			serverBoundKeepalive.PacketID = 0x10
			serverBoundKeepalive.KeepAliveID = receivedKeepalive.KeepAliveID

			err = client.WritePacket(serverBoundKeepalive)
			if err != nil {
				panic(err)
			}
		}
	}
}
