package store

import (
	"fmt"
	"log"

	"github.com/xxx/fangzhang-bot/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(driver, dsn string) *gorm.DB {
	var err error

	// 确保数据目录存在
	if driver == "sqlite3" {
		// 创建data目录
		// 这里可以添加目录创建逻辑
	}

	// 连接数据库
	switch driver {
	case "sqlite3":
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	default:
		log.Fatalf("不支持的数据库驱动: %s", driver)
	}

	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&model.User{},
		&model.Group{},
		&model.Message{},
		&model.Schedule{},
		&model.Lottery{},
		&model.Participant{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("✅ 数据库初始化成功")
	return DB
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
