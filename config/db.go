package config

import (
	"fmt"
	"log"

	"github.com/Notailab/Notailab/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		Cfg.Database.Host,
		Cfg.Database.User,
		Cfg.Database.Password,
		Cfg.Database.DBName,
		Cfg.Database.Port,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("连接数据库失败：%v", err)
		return
	}

	if err := DB.AutoMigrate(
		&model.User{},
		&model.UserSetting{},
		&model.Project{},
		&model.File{},
		&model.AgentConversation{},
		&model.AgentMessage{},
	); err != nil {
		log.Printf("自动迁移数据库失败：%v", err)
	}
}
