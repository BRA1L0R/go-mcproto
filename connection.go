package mcprot

type dataSerialization interface {
	SerializeData(inter interface{}) error
}

type PacketOp interface {
	dataSerialization
	Serialize(compressionTreshold int) ([]byte, error)
}

func (mc *McProt) WritePacket(packet PacketOp) error {
	if err := packet.SerializeData(packet); err != nil {
		return err
	}

	serialized, err := packet.Serialize(mc.compressionTreshold)
	if err != nil {
		return err
	}

	_, err = mc.connection.Write(serialized)
	return err
}

type UncompressedPacketOp interface {
	dataSerialization
	Serialize() []byte
}

func (mc *McProt) WriteUncompressedPacket(packet UncompressedPacketOp) error {
	if err := packet.SerializeData(packet); err != nil {
		return err
	}

	_, err := mc.connection.Write(packet.Serialize())
	return err
}
