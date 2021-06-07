package main

import (
	"flag"

	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/packets"
	"github.com/BRA1L0R/go-mcproto/packets/models"
)

type EntityPositionUpdate struct {
	packets.MinecraftPacket

	EntityID int32 `mc:"varint"`
	DeltaX   int16 `mc:"inherit"`
	DeltaY   int16 `mc:"inherit"`
	DeltaZ   int16 `mc:"inherit"`
	OnGround bool  `mc:"inherit"`
}

type EntityMultiPositionUpdate struct {
	packets.MinecraftPacket

	EntityID int32 `mc:"varint"`
	DeltaX   int16 `mc:"inherit"`
	DeltaY   int16 `mc:"inherit"`
	DeltaZ   int16 `mc:"inherit"`
	Yaw      uint8 `mc:"inherit"`
	Pitch    uint8 `mc:"inherit"`
	OnGround bool  `mc:"inherit"`
}

type PlayerPositionAndLook struct {
	packets.MinecraftPacket

	X          float64 `mc:"inherit"`
	Y          float64 `mc:"inherit"`
	Z          float64 `mc:"inherit"`
	Yaw        float32 `mc:"inherit"`
	Pitch      float32 `mc:"inherit"`
	Flags      byte    `mc:"inherit"`
	TeleportID int32   `mc:"varint"`
}

type TeleportConfirm struct {
	packets.MinecraftPacket

	TeleportID int32 `mc:"varint"`
}

type PlayerPosition struct {
	packets.MinecraftPacket

	X        float64 `mc:"inherit"`
	FeetY    float64 `mc:"inherit"`
	Z        float64 `mc:"inherit"`
	OnGround bool    `mc:"inherit"`
}

type SpawnPlayer struct {
	packets.MinecraftPacket

	EntityID    int32      `mc:"varint"`
	PlayerUUID  complex128 `mc:"inherit"`
	PlayerX     float64    `mc:"inherit"`
	PlayerY     float64    `mc:"inherit"`
	PlayerZ     float64    `mc:"inherit"`
	PlayerYaw   uint8      `mc:"inherit"`
	PlayerPitch uint8      `mc:"inherit"`
}

var (
	host = flag.String("host", "127.0.0.1", "Server host")
	port = flag.Uint("port", 25565, "Server port")
)

func main() {
	flag.Parse()

	client := mcproto.Client{
		Host:            *host,
		Port:            uint16(*port),
		ProtocolVersion: 754, // 1.16.5
		Name:            "FollowBot",
	}

	err := client.Initialize()
	if err != nil {
		panic(err)
	}

	var lastSpawnedPlayer SpawnPlayer

	UpdatePlayerPos := func(DeltaX int16, DeltaY int16, DeltaZ int16) {
		lastSpawnedPlayer.PlayerX += calculateDelta(DeltaX)
		lastSpawnedPlayer.PlayerY += calculateDelta(DeltaY)
		lastSpawnedPlayer.PlayerZ += calculateDelta(DeltaZ)
	}

	FollowPlayer := func() {
		playerPosPacket := new(PlayerPosition)
		playerPosPacket.PacketID = 0x12

		playerPosPacket.X = lastSpawnedPlayer.PlayerX
		playerPosPacket.FeetY = lastSpawnedPlayer.PlayerY
		playerPosPacket.Z = lastSpawnedPlayer.PlayerZ
		playerPosPacket.OnGround = true

		err := client.WritePacket(playerPosPacket)
		if err != nil {
			panic(err)
		}
	}

	for {
		packet, err := client.ReceivePacket()
		if err != nil {
			panic(err)
		}

		switch packet.PacketID {
		case 0x34:
			receivedPositionAndLook := new(PlayerPositionAndLook)
			err := packet.DeserializeData(receivedPositionAndLook)
			if err != nil {
				panic(err)
			}

			teleportConfirm := new(TeleportConfirm)
			teleportConfirm.PacketID = 0x00
			teleportConfirm.TeleportID = receivedPositionAndLook.TeleportID

			err = client.WritePacket(teleportConfirm)
			if err != nil {
				panic(err)
			}
		case 0x1F: // keepalive
			// receivedKeepalive := models.KeepAlivePacket{MinecraftPacket: packet}
			receivedKeepalive := new(models.KeepAlivePacket)
			err := packet.DeserializeData(receivedKeepalive)
			if err != nil {
				panic(err)
			}

			serverBoundKeepalive := new(models.KeepAlivePacket)
			serverBoundKeepalive.PacketID = 0x10
			serverBoundKeepalive.KeepAliveID = receivedKeepalive.KeepAliveID

			err = client.WritePacket(serverBoundKeepalive)
			if err != nil {
				panic(err)
			}
		case 0x04: // spawn player
			err := lastSpawnedPlayer.DeserializeData(&lastSpawnedPlayer)
			if err != nil {
				panic(err)
			}
		case 0x27: // entity position
			// receivedPositionUpdate := EntityPositionUpdate{MinecraftPacket: packet}
			receivedPositionUpdate := new(EntityPositionUpdate)
			err := packet.DeserializeData(receivedPositionUpdate)
			if err != nil {
				panic(err)
			}

			if lastSpawnedPlayer.EntityID != receivedPositionUpdate.EntityID {
				continue
			}

			UpdatePlayerPos(receivedPositionUpdate.DeltaX, receivedPositionUpdate.DeltaY, receivedPositionUpdate.DeltaZ)
			FollowPlayer()
		case 0x28: // entity position and rotation
			// receivedMultiPositionUpdate := EntityMultiPositionUpdate{MinecraftPacket: packet}
			receivedMultiPositionUpdate := new(EntityMultiPositionUpdate)
			err := packet.DeserializeData(receivedMultiPositionUpdate)
			if err != nil {
				panic(err)
			}

			if lastSpawnedPlayer.EntityID != receivedMultiPositionUpdate.EntityID {
				continue
			}

			UpdatePlayerPos(receivedMultiPositionUpdate.DeltaX, receivedMultiPositionUpdate.DeltaY, receivedMultiPositionUpdate.DeltaZ)
			FollowPlayer()
		}
	}
}

func calculateDelta(delta int16) float64 {
	return (float64(delta) / (32 * 128))
}
