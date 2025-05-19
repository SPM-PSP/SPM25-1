package handler

import (
	"UnoBackend/internal/model/Uno"
	"UnoBackend/internal/model/deepseek"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"strings"
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

func ValidateCardPlayHandler(c *gin.Context) {
	type Request struct {
		RoomID      string   `json:"room_id"`
		PlayerIndex int      `json:"player_index"`
		Card        Uno.Card `json:"card"`
	}
	type Response struct {
		Valid bool `json:"valid"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的 request"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到房间"})
		return
	}

	valid := service.ValidateCardPlay(room, req.PlayerIndex, req.Card)
	service.BroadcastToClients(req.RoomID)
	c.JSON(http.StatusOK, Response{Valid: valid})
}

func HandleSpecialCardHandler(c *gin.Context) {
	type Request struct {
		RoomID string   `json:"room_id"`
		Card   Uno.Card `json:"card"`
	}
	type Response struct {
		Message string `json:"message"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	service.HandleSpecialCard(room, req.Card)
	service.BroadcastToClients(req.RoomID)
	c.JSON(http.StatusOK, Response{Message: "Card handled successfully"})
}

func HandleAcceptHandler(c *gin.Context) {
	type Request struct {
		RoomID string `json:"room_id"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	room, _ := service.GetRoom(req.RoomID)
	if room.Direction == Uno.Clockwise {
		room.CurrentPlayerIndex = (room.CurrentPlayerIndex + 1) % len(room.Players)
	}
	if room.Direction == Uno.Anticlockwise {
		room.CurrentPlayerIndex = (room.CurrentPlayerIndex + len(room.Players) - 1) % len(room.Players)
	}
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	service.HandleAcceptCard(room)
	c.JSON(http.StatusOK, room)
}

func DrawCardHandler(c *gin.Context) {
	type Request struct {
		RoomID   string `json:"room_id"`
		PlayerID string `json:"player_id"`
		Number   int    `json:"number"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	room, _ := service.GetRoom(req.RoomID)
	for _, player := range room.Players {
		if player.ID == req.PlayerID {
			err := service.DrawCards(player, req.Number, room)
			if room.Direction == Uno.Clockwise {
				room.CurrentPlayerIndex = (room.CurrentPlayerIndex + 1) % len(room.Players)
			}
			if room.Direction == Uno.Anticlockwise {
				room.CurrentPlayerIndex = (room.CurrentPlayerIndex + len(room.Players) - 1) % len(room.Players)
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "摸牌发生错误"})
				return
			}
			break
		}
	}
	c.JSON(http.StatusOK, room)
}

func UnoChatHandler(c *gin.Context) {
	type Request struct {
		Content       string `json:"content"`
		AiPlayerIndex int    `json:"ai_player_index"`
		RoomId        string `json:"room_id"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}
	room, _ := service.GetRoom(req.RoomId)
	handJson, _ := json.MarshalIndent(room.Players[req.AiPlayerIndex].Hand, "", "    ")
	cardJson, _ := json.MarshalIndent(room.DiscardPile[len(room.DiscardPile)-1], "", "    ")
	message := fmt.Sprintf(
		"现在你正在玩Uno游戏，允许+2后叠加+4，+4后依然可以无限制的叠加+4，若出wild和+4需要返回你所选定的color信息，如果惩罚值为0你可以仅仅凭借颜色进行出牌,如果惩罚值不为0，则只能出+4，不允许抢出，你的手牌是\"hand\": %s 现在场上的牌是 %s 惩罚值为 %d 你需要以JSON格式返回你要打出的牌,且仅仅返回json,如果你无法出牌，则返回空白的json",
		string(handJson), string(cardJson), room.DrawCount,
	)
	fmt.Println(message)
	messages := []deepseek.ChatMessage{
		{Role: "user", Content: fmt.Sprintf("%s", message)},
	}
	//"你好！现在我需要你扮演猫娘来和我进行对话，具体表现为句末带上‘喵～’字样并且语言风格偏向可爱。"
	response, err := service.GetDeepSeekChatCompletion(messages)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	cleaned := strings.TrimPrefix(response, "```json\n")
	cleaned = strings.TrimSuffix(cleaned, "\n```")
	var AiCard Uno.Card
	err1 := json.Unmarshal([]byte(cleaned), &AiCard)
	if err1 != nil {
		panic(err1)
	}
	c.JSON(200, AiCard)
	fmt.Println("Assistant:", response)
}
