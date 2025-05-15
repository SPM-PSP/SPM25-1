package handler

import (
	"UnoBackend/DB"
	"UnoBackend/config"
	"UnoBackend/internal/model/suop"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetAllSuops(c *gin.Context) {
	var suops []suop.Suop
	if err := DB.DB.Find(&suops).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, suops)
}

func CreateSuop(c *gin.Context) {
	var newSuop suop.Suop
	if err := c.ShouldBindJSON(&newSuop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	newSuop.CreatedAt = formattedTime
	if err := DB.DB.Create(&newSuop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newSuop)
}

func UpdateSuop(c *gin.Context) {
	id := c.Param("id")
	var suopData suop.Suop
	if err := DB.DB.First(&suopData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录未找到"})
		return
	}

	if err := c.ShouldBindJSON(&suopData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.DB.Save(&suopData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suopData)
}

// 删除 Suop
func DeleteSuop(c *gin.Context) {
	id := c.Param("id")
	if err := DB.DB.Delete(&suop.Suop{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetSuop(c *gin.Context) {
	type Request struct {
		sid int `json:"sid"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}
	var suopData suop.Suop
	if err := DB.DB.Find(&suopData, req.sid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "汤面未找到"})
		return
	}
	c.JSON(http.StatusOK, suopData)
}

func ContinueChat(c *gin.Context) {
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
	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}

	// 获取已有 session
	cfg := config.Load()
	handler := service.NewChatHandler("sk-09e51faee39f4a9a9358dbd732868b1f", cfg.APITimeout)

	// 这里从 handler.store.sessions 中找回旧 session
	session := handler.GetSessionByID(room.Session)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话未找到"})
		return
	}

	// 发送消息
	reply, err := handler.SendAMessage(session, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 回复失败"})
		return
	}

	room.Message = reply
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}

func StartSuop(c *gin.Context) {
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
	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间未找到"})
		return
	}
	if err := DB.DB.Find(&suopData, req.SuopID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "汤面未找到"})
		return
	}
	cfg := config.Load()
	service.StartSuopGame(room, req.SuopID, service.NewChatHandler("sk-09e51faee39f4a9a9358dbd732868b1f",
		cfg.APITimeout))
	c.JSON(http.StatusOK, room)
}
