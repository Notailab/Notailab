package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/config"
	"github.com/Notailab/Notailab/internal/router"
)

func main() {
	config.InitConfig()
	config.InitDB()

	router := router.InitRouter()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "user.html", gin.H{})
	})

	router.Run()
}