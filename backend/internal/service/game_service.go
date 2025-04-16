package service

import (
	"UnoBackend/internal/model/Uno"
	"math/rand"
	"time"
)

func ValidateCardPlay(room *Uno.Room, playerIndex int, card Uno.Card) bool {
	topCard := room.DiscardPile[len(room.DiscardPile)-1]

	// 万能牌始终合法
	if card.Type == "wild" || card.Type == "wild_draw_four" {
		return true
	}

	// 颜色或数值匹配
	return card.Color == topCard.Color || card.Value == topCard.Value
}

func HandleSpecialCard(room *Uno.Room, card Uno.Card, choose string) {
	switch card.Type {
	case "reverse":
		reversePlayerOrder(room)
		room.DiscardPile = append(room.DiscardPile, card)
	case "draw_two", "draw_four":
		//计算抽牌累计
		if card.Type == "draw_two" {
			room.DrawCount += 2
		} else if card.Type == "draw_four" {
			room.DrawCount += 4
		}
		room.DiscardPile = append(room.DiscardPile, card)
		//自己接受摸牌
		if choose == "accept" {
			err := drawCards(room.Players[room.CurrentPlayerIndex], room.DrawCount, room)
			room.DrawCount = 0
			if err != nil {
				return
			}
		}
	}
}

// 摸牌逻辑
func drawCards(player *Uno.Player, num int, room *Uno.Room) error {
	if len(room.Deck) < num {
		reshuffleDiscardPile(room)
	}

	cards := room.Deck[len(room.Deck)-num:]
	player.Hand = append(player.Hand, cards...)
	room.Deck = room.Deck[:len(room.Deck)-num]

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

func absInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
