package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})

	r := gin.Default()

	// 认证路由组
	auth := r.Group("/auth")
	{
		auth.POST("/register", registerHandler(db))
		auth.POST("/login", loginHandler(db))
	}

	// WebSocket 路由
	r.GET("/ws", WebSocketHandler)

	// 需要认证的路由组
	api := r.Group("/api")
	api.Use(JWTMiddleware())
	{
		api.GET("/user/me", getUserHandler)
	}

	r.Run(":8080")
}
