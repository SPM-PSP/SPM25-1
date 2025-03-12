package model

import (
	"github.com/google/uuid"
	"time"
)

func NewSession() *ChatSession {
	return &ChatSession{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}
}
