package routes

import (
	"UnoBackend/internal/handler"
	"UnoBackend/internal/middle"
	"UnoBackend/internal/service"
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
	router.POST("/Uno/chat", handler.UnoChatHandler)
}

func JoinRoomRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	api := router.Group("/ws")
	{
		api.POST("/joinRoom", handler.JoinRoomHandler)
		api.GET("/joinRoom", handler.WsHandler)
		api.POST("/leaveRoom", handler.LeaveRoomHandler)
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
		api.POST("/getRoomById", handler.GetRoomByIdRoomHandler)
		api.GET("/getUserByUsername", handler.GetUserByUsername)
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

func ValidateCardPlayRoutes(router *gin.Engine) {
	router.Use(middle.CORS())
	api := router.Group("/Uno")
	{
		api.POST("/checkCard", handler.ValidateCardPlayHandler)
		api.POST("/handleSpecial", handler.HandleSpecialCardHandler)
		api.POST("/draw", handler.DrawCardHandler)
		api.POST("/accept", handler.HandleAcceptHandler)
	}
}
