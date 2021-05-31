package packets

type SetCompressionPacket struct {
	*UncompressedPacket

	Treshold int `type:"varint"`
}
