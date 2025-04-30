package handler

import (
	"UnoBackend/internal/model/Uno"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRoomHandler(c *gin.Context) {
	type createRequest struct {
		RoomID  string `json:"room_id"`
		Creator string `json:"creator"`
	}
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	newRoom := service.CreateRoom(req.Creator, req.RoomID)
	c.JSON(http.StatusOK, gin.H{
		"roomID": req.RoomID,
		"room":   newRoom,
	})
}

func GetRoomByIdRoomHandler(c *gin.Context) {
	type getRequest struct {
		RoomID string `json:"room_id"`
	}
	var req getRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}
	fmt.Println("wwww" + req.RoomID)
	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}
	c.JSON(http.StatusOK, room)

}

func JoinRoomHandler(c *gin.Context) {
	type JoinRequest struct {
		RoomID string     `json:"room_id"`
		Player Uno.Player `json:"player"`
	}

	var req JoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	fmt.Println(room.Players)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}

	room.Players = append(room.Players, &req.Player)

	c.JSON(http.StatusOK, gin.H{
		"player": req.Player.ID,
		"room":   room.ID,
	})
}
