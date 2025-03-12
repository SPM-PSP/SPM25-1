package main

import (
	"UnoBackend/DB"
	"UnoBackend/config"
	"UnoBackend/internal/domain"
	"UnoBackend/internal/handler"
	"UnoBackend/internal/middle"
	"UnoBackend/internal/model"
	"UnoBackend/internal/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	DB.InitDB()
	// 用户认证相关
	r.POST("/register", routes.Register)
	r.POST("/login", routes.Login)

	// WebSocket
	r.GET("/ws", routes.WebSocketHandler)

	// 需要 JWT 保护的接口
	auth := r.Group("/")
	auth.Use(middle.JWTAuth())
	auth.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "访问成功"})
	})
	auth.POST("/protected/chat", func(c *gin.Context) {
		messages := []model.ChatMessage{
			{Role: "user", Content: "你好！现在我需要你扮演猫娘来和我进行对话，具体表现为句末带上‘喵～’字样并且语言风格偏向可爱。"},
		}

		response, err := domain.GetDeepSeekChatCompletion(messages)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		c.JSON(200, gin.H{"message": response})
		fmt.Println("Assistant:", response)
	})
	cfg := config.Load()
	fmt.Println("Assistant:", cfg.DeepSeekAPIKey)
	chatHandler := handler.NewChatHandler(
		"sk-09e51faee39f4a9a9358dbd732868b1f",
		cfg.APITimeout,
	)
	routes.RegisterRoutes(r, chatHandler)

	r.Run(":8080")
}
