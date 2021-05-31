package packets

type HandshakePacket struct {
	*UncompressedPacket

	ProtocolVersion int    `type:"varint"`
	ServerAddress   string `type:"string"`
	ServerPort      uint16 `type:"inherit"`
	NextState       int    `type:"varint"`
}
