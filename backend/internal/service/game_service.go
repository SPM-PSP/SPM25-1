package service

import (
	"UnoBackend/DB"
	"UnoBackend/internal/model/Uno"
	"UnoBackend/internal/model/suop"
	"fmt"
	"math/rand"
	"time"
)

func ValidateCardPlay(room *Uno.Room, playerIndex int, card Uno.Card) bool {
	topCard := room.DiscardPile[len(room.DiscardPile)-1]
	fmt.Println(card)
	fmt.Println(topCard)
	if room.DrawCount != 0 {
		if card.Type == topCard.Type || card.Type == "wild_draw_four" {
			fmt.Println("有惩罚的通过")
			return true
		}
	}
	if room.DrawCount == 0 {
		fmt.Println("无惩罚")
		// 万能牌始终合法
		if card.Type == "wild" || card.Type == "wild_draw_four" {
			return true
		}
		// 颜色或数值匹配
		return card.Color == topCard.Color || card.Value == topCard.Value
	}
	return false
}

// 自己接受摸牌
func HandleAcceptCard(room *Uno.Room) {
	if room.DrawCount == 0 {
		DrawCards(room.Players[room.CurrentPlayerIndex], 1, room)
	}
	err := DrawCards(room.Players[room.CurrentPlayerIndex], room.DrawCount, room)
	room.DrawCount = 0
	if err != nil {
		return
	}
}

func HandleSpecialCard(room *Uno.Room, card Uno.Card) {
	room.Players[room.CurrentPlayerIndex].Pre = card
	switch card.Type {
	case "reverse":
		reversePlayerOrder(room)
		RemoveHandCard(room, card)
		room.DiscardPile = append(room.DiscardPile, card)
	case "draw_two", "wild_draw_four":
		//计算抽牌累计
		if card.Type == "draw_two" {
			room.DrawCount += 2
		} else if card.Type == "wild_draw_four" {
			room.DrawCount += 4
		}
		RemoveHandCard(room, card)
		room.DiscardPile = append(room.DiscardPile, card)

	case "number":
		RemoveHandCard(room, card)
		room.DiscardPile = append(room.DiscardPile, card)

	case "skip":
		RemoveHandCard(room, card)

		if room.Direction == Uno.Clockwise {
			room.CurrentPlayerIndex = (room.CurrentPlayerIndex + 1) % len(room.Players)
		}
		if room.Direction == Uno.Anticlockwise {
			room.CurrentPlayerIndex = (room.CurrentPlayerIndex + len(room.Players) - 1) % len(room.Players)
		}
		room.DiscardPile = append(room.DiscardPile, card)
	case "wild":
		RemoveHandCard(room, card)
		room.DiscardPile = append(room.DiscardPile, card)
	}

	if room.Direction == Uno.Clockwise {
		room.CurrentPlayerIndex = (room.CurrentPlayerIndex + 1) % len(room.Players)
	}
	if room.Direction == Uno.Anticlockwise {
		room.CurrentPlayerIndex = (room.CurrentPlayerIndex + len(room.Players) - 1) % len(room.Players)
	}
}

// 移除房间内该回合出牌玩家所出手牌
func RemoveHandCard(room *Uno.Room, card Uno.Card) {
	if card.Type == "wild" || card.Type == "wild_draw_four" {
		card.Color = ""
	}
	fmt.Print(card)
	removed := false
	newHands := make([]Uno.Card, 0, len(room.Players[room.CurrentPlayerIndex].Hand))
	for _, p := range room.Players[room.CurrentPlayerIndex].Hand {
		if !removed && p == card {
			removed = true // 只删除第一次匹配到的牌
			continue
		}
		newHands = append(newHands, p)
	}
	room.Players[room.CurrentPlayerIndex].Hand = newHands

	if !removed {
		fmt.Print("error  玩家没有此手牌")
		return
	}
}

// 摸牌逻辑
func DrawCards(player *Uno.Player, num int, room *Uno.Room) error {
	//如果不够摸就重洗
	if len(room.Deck) < num {
		reshuffleDiscardPile(room)
	}

	cards := room.Deck[len(room.Deck)-num:]
	player.Hand = append(player.Hand, cards...)
	room.Deck = room.Deck[:len(room.Deck)-num]
	player.Pre.Type = "drown"
	return nil
}

// 弃牌堆重洗
func reshuffleDiscardPile(room *Uno.Room) {
	// 保留最后一张弃牌作为起点
	newDeck := room.DiscardPile[:len(room.DiscardPile)-1]
	room.Deck = shuffle(newDeck)
	room.DiscardPile = []Uno.Card{room.DiscardPile[len(room.DiscardPile)-1]}
}

// 洗牌
func shuffle(deck []Uno.Card) []Uno.Card {
	rand.Seed(time.Now().UnixNano()) // 设置随机种子
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i] // 交换元素
	})
	return deck
}

// 下家
func getNextPlayer(room *Uno.Room) *Uno.Player {
	if room.Direction == Uno.Clockwise {
		return room.Players[(room.CurrentPlayerIndex+1)%len(room.Players)]
	} else {
		index := room.CurrentPlayerIndex - 1
		if index < 0 {
			index = len(room.Players) - 1
		}
		return room.Players[index]
	}

}

// 反转
func reversePlayerOrder(room *Uno.Room) {
	if room.Direction == Uno.Clockwise {
		room.Direction = Uno.Anticlockwise
	} else {
		room.Direction = Uno.Clockwise
	}
}

func StartUnoGame(room *Uno.Room) {
	for i := range room.Players {
		DrawCards(room.Players[i], 4, room)
	}
	if room.Deck[0].Color == "" {
		room.Deck[0].Color = Uno.Red
	}
	room.DiscardPile = append(room.DiscardPile, room.Deck[0])
	room.Status = Uno.Playing
	room.Direction = Uno.Clockwise

}

func StartSuopGame(room *Uno.Room, id int, handler *ChatHandler) {
	var suopData suop.Suop
	if err := DB.DB.Find(&suopData, id).Error; err != nil {
		fmt.Println("error: 汤面未找到")
		return
	}
	// 创建会话
	session := handler.NewASession()
	// 构造自定义对话内容，例如从 suopData 生成一个问题
	message := fmt.Sprintf("现在你是海龟汤推理游戏的主持人，你要根据我接下来的提问与汤底进行比对，你需要严格执行以下内容，你仅仅可以回答是,不是,不重要,可能，四种回答。只有当我选择复述故事并且已经正确猜到游戏汤底大部分内容后请告诉我汤底故事并且提示我游戏完成，以下是汤底：%s", suopData.Content)

	// 调用对话接口
	response, err := handler.SendAMessage(session, message)
	if err != nil {
		fmt.Println("调用 AI 接口出错:", err)
		return
	}

	fmt.Println("AI 回复：", response)
	room.Status = Uno.Playing
	room.Message = response
	room.Session = session.ID
	// 如果你需要将 response 存储到 room 中，可加上：
	// room.SomeField = response
}

func absInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
