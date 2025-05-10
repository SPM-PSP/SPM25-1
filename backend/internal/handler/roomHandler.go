package handler

import (
	"UnoBackend/internal/model"
	"UnoBackend/internal/model/Uno"
	"UnoBackend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func CreateRoomHandler(c *gin.Context) {
	type createRequest struct {
		RoomID  string     `json:"room_id"`
		Creator Uno.Player `json:"creator"`
	}
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	newRoom := service.CreateRoom(req.Creator.ID, req.RoomID, req.Creator.Avatar)
	//房主的web加入
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsHandler(c *gin.Context) {
	roomId := c.Query("roomId")
	playerId := c.Query("playerId")
	if roomId == "" || playerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 roomId 或 playerId"})
		return
	}
	// 升级 HTTP 为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	// 获取房间，将连接加入
	room, _ := service.GetRoom(roomId)
	if room == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("未找到房间"))
		return
	}
	service.AddClient(room, conn)
	// 广播“有人加入”系统消息
	service.BroadcastMsg(room, model.Message{
		Type: "system",
		Data: fmt.Sprintf("玩家 %s 加入了房间", playerId),
	})

	// 持续读取该连接的消息
	for {
		var msg model.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}
		// 收到消息后直接广播给同房间的所有客户端
		service.BroadcastMsg(room, msg)
	}

	// 连接中断，移除并广播离线消息
	service.RemoveClient(room, conn)
	service.BroadcastMsg(room, model.Message{
		Type: "system",
		Data: fmt.Sprintf("玩家 %s 离开了房间", playerId),
	})
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

func LeaveRoomHandler(c *gin.Context) {
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

	removed := false
	newPlayers := make([]*Uno.Player, 0, len(room.Players))
	for _, p := range room.Players {
		if p.ID != req.Player.ID {
			newPlayers = append(newPlayers, p)
		} else {
			removed = true
		}
	}
	room.Players = newPlayers

	if !removed {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家未在房间中"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "玩家已离开房间",
		"player":  req.Player.ID,
		"room":    room.ID,
	})
}
