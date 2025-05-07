package model

type Message struct {
	Type string      `json:"type"` // e.g. "join", "chat", "play"
	Data interface{} `json:"data"`
}
