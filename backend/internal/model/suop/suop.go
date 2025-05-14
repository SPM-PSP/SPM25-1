package suop

type Suop struct {
	ID          int    `gorm:"column:sid;primaryKey" json:"sid"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	CreatedAt   string `gorm:"column:create_time" json:"create_time"`
}
