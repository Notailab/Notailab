package main

import (
	"github.com/Notailab/Notailab/config"
	"github.com/Notailab/Notailab/internal/router"
)

func main() {
	config.InitConfig()
	config.InitDB()

	router := router.InitRouter()

	router.Run()
}
