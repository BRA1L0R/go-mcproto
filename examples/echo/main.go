package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/BRA1L0R/go-mcprot"
	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/packets/models"
)

type ClientBoundChatMessage struct {
	*packets.CompressedPacket

	JsonData string     `type:"string"`
	Position byte       `type:"inherit"`
	Sender   complex128 `type:"inherit"`
}

type ServerBoundChatMessage struct {
	*packets.CompressedPacket

	Message string `type:"string"`
}

type ChatMessage struct {
	Translate string `json:"translate"`

	With []struct {
		Text string `json:"text"`
	}
}

// var (
// 	host = flag.String("host", "127.0.0.1", "Server host")
// 	port = flag.Uint("port", 25565, "Server port")
// )

func main() {
	// flag.Parse()

	client := mcprot.McProt{
		Host:            "152.89.170.123",
		Port:            25565,
		ProtocolVersion: 754, // 1.16.5
		Name:            "Echo_Bot",
	}

	err := client.Initialize()
	if err != nil {
		panic(err)
	}

	botInfoMessage := ServerBoundChatMessage{
		CompressedPacket: packets.NewCompressedPacket(0x03),
		Message:          fmt.Sprintf("I'm running on %v, %v", runtime.GOOS, runtime.GOARCH),
	}

	err = client.WritePacket(&botInfoMessage)
	if err != nil {
		panic(err)
	}

	for {
		packet, err := client.ReceivePacket()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch packet.PacketID {
		case 0x0E:
			receivedChatMessage := ClientBoundChatMessage{CompressedPacket: packet}

			err := receivedChatMessage.DeserializeData(&receivedChatMessage)
			if err != nil {
				panic(err)
			}

			chatMessage := new(ChatMessage)
			err = json.Unmarshal([]byte(receivedChatMessage.JsonData), chatMessage)
			if err != nil {
				panic(err)
			}

			if chatMessage.Translate == "chat.type.text" {
				user := chatMessage.With[0].Text
				playerText := chatMessage.With[1].Text

				if user == client.Name {
					continue
				}

				fmt.Printf("<%s> %s\n", user, playerText)

				chatMessage := ServerBoundChatMessage{
					CompressedPacket: packets.NewCompressedPacket(0x03),
					Message:          playerText,
				}

				err := client.WritePacket(&chatMessage)
				if err != nil {
					panic(err)
				}
			}
		case 0x1F:
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
		}
	}
}
