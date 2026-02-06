package config

import (
	"fmt"
	"log"

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
	}
}