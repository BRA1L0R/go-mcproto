package mcproto

// PacketOp defines the standard methods that a struct should have
// in order to be serializable by the library
//
// You can actually create your own methods as long as they respect this standard
type PacketOp interface {
	SerializeData(inter interface{}) error
	Serialize(compressionTreshold int32) ([]byte, error)
}

func (mc *McProto) WritePacket(packet PacketOp) error {
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
