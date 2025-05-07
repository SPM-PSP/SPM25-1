package handler

import (
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
