package Uno

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	ID   string          `json:"id"`
	Conn *websocket.Conn `json:"-"`
	Hand []Card          `json:"hand"`
}
