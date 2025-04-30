package suop

type Suop struct {
	ID      int    `gorm:"column:sid;primaryKey" json:"sid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
