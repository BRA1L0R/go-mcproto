package types

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PlayerSample struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`

	Sample []PlayerSample `json:"sample"`
}

type Description struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Version     Version     `json:"version"`
	Players     Players     `json:"players"`
	Description Description `json:"description"`
	Favicon     string      `json:"favicon"`
}
