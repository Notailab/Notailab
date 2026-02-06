package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

var Cfg *Config

func InitConfig()  {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败：%v", err)
	}
	
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("解析配置文件失败：%v", err)
	}
	
	fmt.Println("配置初始化成功！")
}