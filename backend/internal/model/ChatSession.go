package model

import "time"

type ChatSession struct {
	ID         string        `json:"id"`
	Messages   []ChatMessage `json:"messages"`
	CreatedAt  time.Time     `json:"created_at"`
	LastActive time.Time     `json:"last_active"`
}
