<div align="center">
<img src="assets/illustration-transparent.png"/>
  <a href="https://pkg.go.dev/github.com/BRA1L0R/go-mcproto"><img src="https://pkg.go.dev/badge/github.com/BRA1L0R/go-mcproto.svg" alt="Go Reference"></a>
  <a href="http://github.com/BRA1L0R/go-mcproto"><img src="https://img.shields.io/github/go-mod/go-version/BRA1L0R/go-mcproto"></a>
</div>

<div align="right">
<h6> illustration by <a href="github.com/talentlessguy" >@talentlessguy </a></h6>
</div>

## Install

You'll need to have at least Go 16.x+ to use this library

```sh
go get github.com/BRA1L0R/go-mcproto
```

## Opening a connection

```go
client := mcproto.McProto{Host: "IP or Hostname", Port: 25565, Name: "GoBot", ProtocolVersion: 754}
```

Before opening a connection to a server, you'll have to specified some vital information such as the host, the port,
the name of your bot (which will login in offline mode) and the protocol version, which will have to match the server one,
unless the server uses some kind of backwards compatibility plugin such as ViaVersion.

In this case 754 is the protocol version for Minecraft 1.16.5, but you can find all the versions [here](https://wiki.vg/Protocol_version_numbers)

Once you define client, you can either use the already implemented handshake method (`client.Initialize`) or you can manually open the connection using `client.Connect`

`client.Initialize` only supports offline mode (unauthenticated SP) at the moment, but I will soon implement online mode as well.

## Defining a packet

Not all packets are implemented in the library, in fact, there are only the ones that will get you past the login state.

Fortunately, defining a packet with go-mcproto is as easy as declaring a struct. Here's a quick example

```go
type ChatMessage struct {
  packets.MinecraftPacket

  JsonData string `mc:"string"`
  Position byte   `mc:"inherit"`
  Sender   []byte `mc:"bytes" len:"16"` // UUID is 128bit, hence the 16 byte len
}
```

packets.MinecraftPackets will add the rest of the fields which are needed for a packet to be conformant to the standard format, compressed or not

ChatMessage will also inherit from packets.MinecraftPackets the methods for the serialization of the fields into Data and the final serialization which returns the byte slice that can be sent over the connection using `client.WritePacket`

### Array Length

Sometimes you might encounter a packet which sends before an array the length of it. Fortunately you can still deserialize the packet with no extra steps doing something like this:

```go
type Biome struct {
  BiomeID int32 `mc:"varint"`
}

type ChunkData struct {
  packets.MinecraftPacket

  BiomesLength int32   `mc:"varint"`
  Biomes       []Biome `mc:"array" len:"BiomesLength"`
}
```

The Biome struct can contain as much fields as you like, as the minecraft protocol has a lot of cases where there are multi-fielded arrays

### Field Dependency

In some cases, a field is present only if the previous one is true (in the case of a boolean). You can still manage to do this with only struct tags:

```go
type PlayerInfo struct {
  packets.MinecraftPacket

  HasDisplayName bool   `mc:"inherit"`
  DisplayName    string `mc:"string" depends_on:"HasDisplayName"`
}
```

## Sending a packet

Once you have defined a packet you want to send, you will have to call `client.WritePacket()` to serialize the data and send it over the connection.

```go
packet := new(MyPacket)
packet.PacketID = 0x00 // you can look up the packet ids on wiki.vg
// all the other fields

client.WritePacket(packet)
```

## Receiving a packet

TODO

## Example ✨

This example initializes the connection between the client and the server, thus switching to the Play state, and listens for keepalive packets to which it responds

```go
package main

import (
	"github.com/BRA1L0R/go-mcproto"
	"github.com/BRA1L0R/go-mcproto/packets/models"
)

func main() {
  client := mcproto.McProto{
    Host: "my.minecraftserver.com",
    Port: 25565,
    ProtocolVersion: 754, // 1.16.5
    Name: "QuickStart",
  }

  client.Initialize()

  for {
    packet, err := client.ReceivePacket()
    if err != nil {
      panic(err)
    }

    if packet.PacketID == 0x1F { // clientbound keepalive packetid
      receivedKeepalive := models.KeepAlivePacket{MinecraftPacket: packet}

      err := receivedKeepalive.DeserializeData(&receivedKeepalive)
      if err != nil {
        panic(err)
      }

      serverBoundKeepalive := new(models.KeepAlivePacket)
      serverBoundKeepalive.KeepAliveID = receivedKeepalive.KeepAliveID
      serverBoundKeepalive.PacketID = 0x10 // serverbound keepaliveid

      err = client.WritePacket(serverBoundKeepalive)
    }
  }
}
```
