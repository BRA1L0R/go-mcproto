package packets

type LoginStartPacket struct {
	*UncompressedPacket

	Name string `type:"string"`
}
