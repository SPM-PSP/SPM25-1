package service

import (
	"UnoBackend/internal/model"
	"fmt"
	"strconv"
	"sync"
)

var (
	rooms sync.Map // 全局房间存储
)

func CreateRoom(creatorID string) *model.Room {
	roomID := model.NewRoom()
	Waiting := model.Waiting
	newRoom := &model.Room{
		ID:      roomID.ID,
		Players: []*model.Player{{ID: creatorID}},
		Deck:    initializeDeck(),
		Status:  Waiting,
		Creator: creatorID,
	}
	rooms.Store(roomID, newRoom)
	return newRoom
}

func GetRoom(roomID string) (*model.Room, bool) {
	val, ok := rooms.Load(roomID)
	if !ok {
		return nil, false
	}
	return val.(*model.Room), true
}

// 初始化UNO牌堆
func initializeDeck() []model.Card {
	// 实现108张牌的生成逻辑
	var deck []model.Card
	colors := []model.Color{model.Red, model.Blue, model.Green, model.Yellow}

	// 生成彩色卡牌（数字牌和功能牌）
	for _, color := range colors {
		// 数字牌（0-9）
		// 数字0每个颜色1张
		deck = append(deck, model.Card{
			Type:  "number",
			Color: color,
			Value: "0",
		})

		// 数字1-9每个颜色2张
		for num := 1; num <= 9; num++ {
			value := strconv.Itoa(num)
			deck = append(deck, model.Card{
				Type:  "number",
				Color: color,
				Value: value,
			})
			deck = append(deck, model.Card{
				Type:  "number",
				Color: color,
				Value: value,
			})
		}

		// 功能牌（每种2张）
		actions := []struct {
			cardType string
			value    string
		}{
			{"skip", "skip"},
			{"reverse", "reverse"},
			{"draw_two", "draw_two"},
		}

		for _, action := range actions {
			for i := 0; i < 2; i++ {
				deck = append(deck, model.Card{
					Type:  action.cardType,
					Color: color,
					Value: action.value,
				})
			}
		}
	}

	// 生成万能牌（每种4张）
	wildCards := []struct {
		cardType string
		value    string
	}{
		{"wild", "wild"},
		{"wild_draw_four", "wild_draw_four"},
	}

	for _, wild := range wildCards {
		for i := 0; i < 4; i++ {
			deck = append(deck, model.Card{
				Type:  wild.cardType,
				Color: "", // 万能牌没有颜色
				Value: wild.value,
			})
		}
	}
	shuffled := shuffle(deck)
	fmt.Print(shuffled)
	return shuffled
}
