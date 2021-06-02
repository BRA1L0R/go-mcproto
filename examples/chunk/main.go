package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
	"github.com/BRA1L0R/go-mcproto/varint"
)

type HeightMaps struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
}

type BlockEntity struct {
	X int `nbt:"x"`
	Y int `nbt:"y"`
	Z int `nbt:"z"`
}

type ChunkPacket struct {
	*packets.CompressedPacket

	ChunkX         int32      `type:"inherit"`
	ChunkY         int32      `type:"inherit"`
	FullChunk      bool       `type:"inherit"`
	PrimaryBitMask int32      `type:"varint"`
	HeightMaps     HeightMaps `type:"nbt"`

	BiomesArrayLength int32   `type:"varint"    depends_on:"FullChunk"`
	Biomes            []int32 `type:"varintarr" depends_on:"FullChunk" len:"BiomesArrayLength"`

	ChunkDataSize int32  `type:"varint"`
	ChunkData     []byte `type:"bytes" len:"ChunkDataSize"`

	BlockEntitiesArrayLength int32         `type:"varint"`
	BlockEntities            []BlockEntity `type:"nbt" len:"BlockEntitiesArrayLength"`
}

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
		Name:            "GolangComp",
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

		switch packet.PacketID {
		case 0x20:
			chunk := ChunkPacket{CompressedPacket: packet}
			err := chunk.DeserializeData(&chunk)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Chunk X: %v, Chunk Y: %v\nSpecial Block Locations: %v\n", chunk.ChunkX, chunk.ChunkY, chunk.BlockEntities)
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
