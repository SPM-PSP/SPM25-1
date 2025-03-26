package model

type Color string

const (
	Red    Color = "red"
	Blue   Color = "blue"
	Green  Color = "green"
	Yellow Color = "yellow"

	// ...其他颜色
)

type Card struct {
	Type  string `json:"type"` // number, skip, reverse, draw_two, wild, wild_draw_four
	Color Color  `json:"color"`
	Value string `json:"value"` // 0-9, "skip"等
}
