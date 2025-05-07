package routes

import (
	"UnoBackend/internal/handler"
	"UnoBackend/internal/middle"
	"UnoBackend/internal/model/deepseek"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(router *gin.Engine, chatHandler *service.ChatHandler) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/deepseek")
	{
		api.POST("/sessions", chatHandler.CreateSession)
		api.POST("/chat", chatHandler.HandleChat)
	}
}

func RegisterRegisterRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	api := router.Group("/")
	{
		api.POST("/register", handler.Register)
	}
}

func RegisterLoginRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	api := router.Group("/")
	{
		api.POST("/login", handler.Login)
	}
}

func RegisterUnoChatRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "访问成功"})
	})
	router.POST("/Uno/chat", func(c *gin.Context) {
		messages := []deepseek.ChatMessage{
			{Role: "user", Content: "你好！现在我需要你扮演猫娘来和我进行对话，具体表现为句末带上‘喵～’字样并且语言风格偏向可爱。"},
		}

		response, err := service.GetDeepSeekChatCompletion(messages)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		c.JSON(200, gin.H{"message": response})
		fmt.Println("Assistant:", response)
	})
}

func JoinRoomRoutes(router *gin.Engine) {
	api := router.Group("/ws")
	{
		api.POST("/joinRoom", handler.JoinRoomHandler)
		api.GET("/joinRoom", handler.WsHandler)
	}
}

func GetAllSuopsRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/")
	{
		api.GET("/getAllSuop", handler.GetAllSuops)
	}
}

func GetRoomByIdRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/")
	{
		api.GET("/getRoomById", handler.GetRoomByIdRoomHandler)
	}
}

func CreateSuopRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/")
	{
		api.POST("/createSuop", handler.CreateSuop)
	}
}

func CreateRoomRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/")
	{
		api.POST("/createRoom", handler.CreateRoomHandler)
	}
}

func StartUnoRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/")
	{
		api.POST("/StartUno", handler.StartUno)
	}
}

func StartSuopRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	router.Use(middle.JWTAuth())
	api := router.Group("/")
	{
		api.POST("/StartSuop", handler.StartSuop)
	}
}
