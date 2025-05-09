package handler

import (
	"UnoBackend/internal/model/Uno"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartUno(c *gin.Context) {
	type startSRequest struct {
		RoomID string `json:"room_id"`
	}
	var req startSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}
	fmt.Println("ROOMID:" + req.RoomID)
	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}
	service.StartUnoGame(room)
	service.BroadcastToClients(req.RoomID)
	c.JSON(http.StatusOK, room)
}

func ValidateCardPlayHandler(c *gin.Context) {
	type Request struct {
		RoomID      string   `json:"room_id"`
		PlayerIndex int      `json:"player_index"`
		Card        Uno.Card `json:"card"`
	}
	type Response struct {
		Valid bool `json:"valid"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的 request"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到房间"})
		return
	}

	valid := service.ValidateCardPlay(room, req.PlayerIndex, req.Card)
	service.BroadcastToClients(req.RoomID)
	c.JSON(http.StatusOK, Response{Valid: valid})
}

func HandleSpecialCardHandler(c *gin.Context) {
	type Request struct {
		RoomID string   `json:"room_id"`
		Card   Uno.Card `json:"card"`
		Choose string   `json:"choose"` // "accept" 或 "reject"
	}
	type Response struct {
		Message string `json:"message"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	service.HandleSpecialCard(room, req.Card, req.Choose)
	service.BroadcastToClients(req.RoomID)
	c.JSON(http.StatusOK, Response{Message: "Card handled successfully"})
}
func DrawCardHandler(c *gin.Context) {
	type Request struct {
		RoomID   string `json:"room_id"`
		PlayerID string `json:"player_id"`
		Number   int    `json:"number"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	for _, player := range room.Players {
		if player.ID == req.PlayerID {
			err := service.DrawCards(player, req.Number, room)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "摸牌发生错误"})
				return
			}
			break
		}
	}
	c.JSON(http.StatusOK, room)
}
