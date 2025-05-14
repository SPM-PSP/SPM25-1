package service

import (
	"UnoBackend/DB"
	"UnoBackend/internal/model/deepseek"
	"UnoBackend/internal/model/suop"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type SessionStore struct {
	sync.RWMutex
	sessions map[string]*deepseek.ChatSession
}

type ChatHandler struct {
	apiKey     string
	apiTimeout time.Duration
	store      *SessionStore
}

func NewChatHandler(apiKey string, timeout time.Duration) *ChatHandler {
	return &ChatHandler{
		apiKey:     apiKey,
		apiTimeout: timeout,
		store: &SessionStore{
			sessions: make(map[string]*deepseek.ChatSession),
		},
	}
}

func (h *ChatHandler) CreateSession(c *gin.Context) {
	session := deepseek.NewSession()

	h.store.Lock()
	h.store.sessions[session.ID] = session
	h.store.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"session_id": session.ID,
		"created_at": session.CreatedAt.Format(time.RFC3339),
		"api_key":    h.apiKey,
	})

}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id" binding:"required"`
		Message   string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	h.store.RLock()
	session, exists := h.store.sessions[req.SessionID]
	h.store.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	// Add user message
	session.Messages = append(session.Messages, deepseek.ChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	// Call API
	response, err := h.callDeepSeekAPI(session.Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add assistant response
	session.Messages = append(session.Messages, deepseek.ChatMessage{
		Role:    "assistant",
		Content: response,
	})
	session.LastActive = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"response":   response,
		"session_id": session.ID,
	})
}

func (h *ChatHandler) callDeepSeekAPI(messages []deepseek.ChatMessage) (string, error) {
	client := &http.Client{Timeout: h.apiTimeout}

	reqBody := deepseek.ChatCompletionRequest{
		Model:     "deepseek-chat",
		Messages:  messages,
		MaxTokens: 500,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d status", resp.StatusCode)
	}

	var apiResp deepseek.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return apiResp.Choices[0].Message.Content, nil
}

func (h *ChatHandler) NewASession() *deepseek.ChatSession {
	session := deepseek.NewSession()
	h.store.Lock()
	h.store.sessions[session.ID] = session
	h.store.Unlock()
	return session
}

func (h *ChatHandler) SendAMessage(session *deepseek.ChatSession, message string) (string, error) {
	session.Messages = append(session.Messages, deepseek.ChatMessage{
		Role:    "user",
		Content: message,
	})

	response, err := h.callDeepSeekAPI(session.Messages)
	if err != nil {
		return "", err
	}

	session.Messages = append(session.Messages, deepseek.ChatMessage{
		Role:    "assistant",
		Content: response,
	})
	session.LastActive = time.Now()

	return response, nil
}

func (h *ChatHandler) GetSessionByID(id string) *deepseek.ChatSession {
	h.store.Lock()
	defer h.store.Unlock()
	return h.store.sessions[id]
}

func (h *ChatHandler) ContinueChat(c *gin.Context) {
	type ChatRequest struct {
		RoomID  string `json:"room_id"`
		Message string `json:"message"`
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 获取房间
	room, _ := GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}

	// 这里从 handler.store.sessions 中找回旧 session
	session := h.GetSessionByID(room.Session)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话未找到"})
		return
	}

	// 发送消息
	reply, err := h.SendAMessage(session, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 回复失败"})
		return
	}

	room.Message = reply
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}

func (h *ChatHandler) StartSuop(c *gin.Context) {
	type startSRequest struct {
		RoomID  string `json:"room_id"`
		SuopID  int    `json:"suop_id"`
		Session string `json:"session"`
	}
	var req startSRequest
	var suopData suop.Suop
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}
	fmt.Println("ROOMID:" + req.RoomID)
	room, _ := GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}
	if err := DB.DB.Find(&suopData, req.SuopID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "汤面未找到"})
		return
	}
	StartSuopGame(room, req.SuopID, h)
	c.JSON(http.StatusOK, room)
}
