package model

import (
	"github.com/google/uuid"
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
	ID                 string        `json:"id"`
	Players            []*Player     `json:"players"`
	Deck               []Card        `json:"-"`
	DiscardPile        []Card        `json:"discardPile"`
	CurrentPlayerIndex int           `json:"currentPlayerIndex"`
	Status             RoomStatus    `json:"status"`
	Creator            string        `json:"creator"` // 房主ID
	Direction          roomDirection `json:"direction"`
}

func NewRoom() *Room {
	return &Room{
		ID: uuid.New().String(),
	}
}
