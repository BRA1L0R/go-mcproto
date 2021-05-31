package packets

type KeepAlivePacket struct {
	*StandardPacket

	KeepAliveID int64 `type:"inherit"`
}
