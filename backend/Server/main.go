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

	r := gin.Default()

	DB.InitDB()
	// 用户认证相关
	routes.RegisterRegisterRoutes(r)
	routes.RegisterLoginRoutes(r)
	routes.GetAllSuopsRoutes(r)
	routes.CreateSuopRoutes(r)
	routes.CreateRoomRoutes(r)
	routes.JoinRoomRoutes(r)
	routes.GetRoomByIdRoutes(r)
	routes.StartSuopRoutes(r)
	routes.StartUnoRoutes(r)
	//test
	//newRoom := service.CreateRoom("niumo")
	//fmt.Println("New room:", newRoom)
	// WebSocket
	//r.GET("/ws", routes.WebSocketHandler)

	// 需要 JWT 保护的接口
	routes.RegisterUnoChatRoutes(r)
	routes.RegisterChatRoutes(r, chatHandler)
	r.Run(":8080")
}
