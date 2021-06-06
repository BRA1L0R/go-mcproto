<div align="center">
<img src="assets/illustration-transparent.png"/>
  <a href="https://pkg.go.dev/github.com/BRA1L0R/go-mcproto"><img src="https://pkg.go.dev/badge/github.com/BRA1L0R/go-mcproto.svg" alt="Go Reference"></a>
  <a href="http://github.com/BRA1L0R/go-mcproto"><img src="https://img.shields.io/github/go-mod/go-version/BRA1L0R/go-mcproto"></a>
</div>

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
