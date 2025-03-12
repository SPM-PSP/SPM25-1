package DB

import (
	"UnoBackend/internal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:qweasd@tcp(127.0.0.1:3306)/table_game?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("数据表迁移失败:", err)
	}

	fmt.Println("数据库初始化完成")
}
