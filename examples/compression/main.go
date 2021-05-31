package main

import (
	"flag"
	"fmt"

	"github.com/BRA1L0R/go-mcprot"
	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/varint"
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
	*packets.StandardPacket

	ChunkX         int32      `type:"inherit"`
	ChunkY         int32      `type:"inherit"`
	FullChunk      bool       `type:"inherit"`
	PrimaryBitMask int        `type:"varint"`
	HeightMaps     HeightMaps `type:"nbt"`

	BiomesArrayLength int   `type:"varint"    depends_on:"FullChunk"`
	Biomes            []int `type:"varintarr" depends_on:"FullChunk" len:"BiomesArrayLength"`

	ChunkDataSize int    `type:"varint"`
	ChunkData     []byte `type:"bytes" len:"ChunkDataSize"`

	BlockEntitiesArrayLength int           `type:"varint"`
	BlockEntities            []BlockEntity `type:"nbt" len:"BlockEntitiesArrayLength"`
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
		Name:            "GolangComp",
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

		switch packet.PacketID {
		case 0x20:
			chunk := ChunkPacket{StandardPacket: packet}
			err := chunk.DeserializeData(&chunk)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Chunk X: %v, Chunk Y: %v\nSpecial Block Locations: %v\n", chunk.ChunkX, chunk.ChunkY, chunk.BlockEntities)
		case 0x1F:
			receivedKeepalive := packets.KeepAlivePacket{StandardPacket: packet}

			err := receivedKeepalive.DeserializeData(&receivedKeepalive)
			if err != nil {
				panic(err)
			}

			serverBoundKeepalive := packets.KeepAlivePacket{
				StandardPacket: &packets.StandardPacket{PacketID: 0x10},
				KeepAliveID:    receivedKeepalive.KeepAliveID,
			}

			client.WritePacket(serverBoundKeepalive)
		}
	}
}
