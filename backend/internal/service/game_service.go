package service

import (
	"UnoBackend/internal/model"
	"math/rand"
	"time"
)

func ValidateCardPlay(room *model.Room, playerIndex int, card model.Card) bool {
	topCard := room.DiscardPile[len(room.DiscardPile)-1]

	// 万能牌始终合法
	if card.Type == "wild" || card.Type == "wild_draw_four" {
		return true
	}

	// 颜色或数值匹配
	return card.Color == topCard.Color || card.Value == topCard.Value
}

func HandleSpecialCard(room *model.Room, card model.Card) {
	switch card.Type {
	case "reverse":
		reversePlayerOrder(room)
	case "draw_two":
		nextPlayer := getNextPlayer(room)
		//下家有+2、+4

		//下家接受摸牌
		err := drawCards(nextPlayer, 2, room)
		if err != nil {
			return
		}
		// ...处理其他特殊牌
	case "draw_four":
		nextPlayer := getNextPlayer(room)
		//下家有+4

		//下家选择摸牌
		err := drawCards(nextPlayer, 4, room)
		if err != nil {
			return
		}
	}
}

// 摸牌逻辑
func drawCards(player *model.Player, num int, room *model.Room) error {
	if len(room.Deck) < num {
		reshuffleDiscardPile(room)
	}

	cards := room.Deck[len(room.Deck)-num:]
	player.Hand = append(player.Hand, cards...)
	room.Deck = room.Deck[:len(room.Deck)-num]

	return nil
}

// 弃牌堆重洗
func reshuffleDiscardPile(room *model.Room) {
	// 保留最后一张弃牌作为起点
	newDeck := room.DiscardPile[:len(room.DiscardPile)-1]
	room.Deck = shuffle(newDeck)
	room.DiscardPile = []model.Card{room.DiscardPile[len(room.DiscardPile)-1]}
}

// 洗牌
func shuffle(deck []model.Card) []model.Card {
	rand.Seed(time.Now().UnixNano()) // 设置随机种子
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i] // 交换元素
	})
	return deck
}

// 下家
func getNextPlayer(room *model.Room) *model.Player {
	if room.Direction == model.Clockwise {
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
func reversePlayerOrder(room *model.Room) {
	if room.Direction == model.Clockwise {
		room.Direction = model.Anticlockwise
	} else {
		room.Direction = model.Clockwise
	}
}

func absInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
