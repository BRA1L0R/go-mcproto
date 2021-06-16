package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/examples/server/types"
	"github.com/BRA1L0R/go-mcproto/packets"
)

type Handshake struct {
	packets.MinecraftPacket

	ProtocolVersion int32  `mc:"varint"`
	ServerAddress   string `mc:"string"`
	ServerPort      uint16 `mc:"inherit"`
	NextState       int32  `mc:"varint"`
}

type Disconnect struct {
	packets.MinecraftPacket

	Reason string `mc:"string"`
}

type Ping struct {
	packets.MinecraftPacket

	Payload int64 `mc:"inherit"`
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleClient(client *mcproto.Client) {
	// ignore handshake packet
	handshakePacket, err := client.ReceivePacket()
	panicErr(err)

	handshake := new(Handshake)
	panicErr(handshakePacket.DeserializeData(handshake))

	if handshake.NextState != 0x01 {
		disconnect := new(Disconnect)
		disconnect.PacketID = 0x00
		disconnect.Reason = "\"Only ping\""

		panicErr(client.WritePacket(disconnect))
		return
	}

	fmt.Printf("New ping from %s!\n", client.Host)

	for {
		statusPacket, err := client.ReceivePacket()
		if err == io.EOF {
			break
		}

		panicErr(err)

		switch statusPacket.PacketID {
		case 0x00:
			favicon, err := os.ReadFile("gopher.png")
			panicErr(err)

			status := types.StatusResponse{
				Version: types.Version{
					Name:     "1.17",
					Protocol: 755,
				},
				Players: types.Players{
					Max:    -1,
					Online: 2,

					Sample: []types.PlayerSample{
						{Name: "GolangBot1", Id: "4566e69f-c907-48ee-8d71-d7ba5aa00d20"},
						{Name: "AwesomePlayer2", Id: "4566e69f-c907-48ee-8d71-d7ba5aa00d30"},
					},
				},
				Description: types.Description{
					Text: "Awesome minecraft server! Made with go-mcproto!",
				},
				Favicon: fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(favicon)),
			}

			disconnect := new(Disconnect)
			disconnect.PacketID = 0x00

			statusMarshaled, err := json.Marshal(status)
			panicErr(err)

			disconnect.Reason = string(statusMarshaled)

			panicErr(client.WritePacket(disconnect))
		case 0x01:
			ping := new(Ping)

			panicErr(statusPacket.DeserializeData(ping))

			ping.PacketID = 0x01
			panicErr(client.WritePacket(ping))
			panicErr(client.CloseConnection())
			return
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":25565")
	if err != nil {
		panic(err)
	}

	for {
		client, err := mcproto.FromListener(ln)
		if err != nil {
			fmt.Println(err)
			continue
		}

		go HandleClient(client)
	}
}
