package packets

type LoginSuccessPacket struct {
	*StandardPacket

	UUID     string `type:"ignore" len:"16"`
	Username string `type:"string"`
}
