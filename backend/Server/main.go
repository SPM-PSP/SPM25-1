package main

import (
	"UnoBackend/internal/middle"
	"UnoBackend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	database.InitDB()

	// 用户认证相关
	r.POST("/register", routes.Register)
	r.POST("/login", routes.Login)

	// WebSocket
	r.GET("/ws", routes.WebSocketHandler)

	// 需要 JWT 保护的接口
	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())
	auth.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "访问成功"})
	})

	r.Run(":8080")
}
