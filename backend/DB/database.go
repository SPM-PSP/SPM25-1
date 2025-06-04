package DB

import (
	"UnoBackend/config"
	"UnoBackend/internal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dir, _ := os.Getwd()
	fmt.Println("当前工作目录：", dir)
	config.LoadConfig("config.yaml") // 从当前目录读取配置
	dbConfig := config.AppConfig1.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName)
	//dsn := "root:qweasd@tcp(119.84.246.217:38158)/table_game?charset=utf8mb4&parseTime=True&loc=Local"
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
