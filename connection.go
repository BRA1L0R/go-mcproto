package mcproto

type PacketOp interface {
	SerializeData(inter interface{}) error
	Serialize(compressionTreshold int32) ([]byte, error)
	InitializePacket()
}

func (mc *McProto) WritePacket(packet PacketOp) error {
	packet.InitializePacket()

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
