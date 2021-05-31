package packets

type ChatPacket struct {
	*StandardPacket

	JsonData string      `type:"string"`
	Position byte        `type:"byte"`
	Sender   interface{} `type:"ignore" len:"16"`
}
