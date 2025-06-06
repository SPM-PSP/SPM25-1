package main

import (
	"UnoBackend/DB"
	"UnoBackend/config"
	"UnoBackend/internal/routes"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	fmt.Println("Assistant:", cfg.DeepSeekAPIKey)
	chatHandler := service.NewChatHandler(
		"sk-09e51faee39f4a9a9358dbd732868b1f",
		cfg.APITimeout,
	)
	config.LoadConfig("config.yaml") // 从当前目录读取配置
	serverConfig := config.AppConfig1.Server

	r := gin.Default()

	DB.InitDB()
	// 用户认证相关
	routes.JoinRoomRoutes(r)
	routes.RegisterRegisterRoutes(r)
	routes.RegisterLoginRoutes(r)
	routes.GetAllSuopsRoutes(r)
	routes.CreateSuopRoutes(r)
	routes.CreateRoomRoutes(r)
	routes.GetRoomByIdRoutes(r)
	routes.StartSuopRoutes(r, chatHandler)
	routes.StartUnoRoutes(r)
	routes.ValidateCardPlayRoutes(r)
	//test
	//newRoom := service.CreateRoom("niumo")
	//fmt.Println("New room:", newRoom)
	// WebSocket
	//r.GET("/ws", routes.WebSocketHandler)

	// 需要 JWT 保护的接口
	routes.RegisterUnoChatRoutes(r)
	routes.RegisterChatRoutes(r, chatHandler)
	r.Run(":" + serverConfig.Port)
}
