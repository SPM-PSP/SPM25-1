package Uno

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type RoomStatus string
type roomDirection string

const (
	Waiting       RoomStatus    = "waiting"
	Playing       RoomStatus    = "playing"
	Clockwise     roomDirection = "clockwise"
	Anticlockwise roomDirection = "anticlockwise"
)

type Room struct {
	ID                 string                   `json:"id"`
	Players            []*Player                `json:"players"`
	Deck               []Card                   `json:"-"`
	DiscardPile        []Card                   `json:"discardPile"`
	CurrentPlayerIndex int                      `json:"currentPlayerIndex"`
	Status             RoomStatus               `json:"status"`
	Creator            string                   `json:"creator"` // 房主ID
	Direction          roomDirection            `json:"direction"`
	DrawCount          int                      `json:"drawCount"`
	Message            string                   `json:"message"`
	Session            string                   `json:"session"`
	Clients            map[*websocket.Conn]bool `json:"-"`
	Mutex              sync.Mutex               `json:"-"`
}

func NewRoom() *Room {
	return &Room{
		ID: uuid.New().String(),
	}
}
