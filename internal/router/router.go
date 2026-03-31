package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/config"
	"github.com/Notailab/Notailab/internal/dao"
	v1 "github.com/Notailab/Notailab/internal/router/api/v1"
	"github.com/Notailab/Notailab/internal/service"
)

func initStatic(router *gin.Engine) {
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
}

func initUser(api *gin.RouterGroup) {
	userDAO := dao.NewUserDAO(config.DB)
	userService := service.NewUserService(userDAO)
	userHandler := v1.NewUserHandler(userService)
	userHandler.LoadRouter(api)
}

func initProject(api *gin.RouterGroup) {
	projectDAO := dao.NewProjectDAO(config.DB)
	projectService := service.NewProjectService(projectDAO)
	projectHandler := v1.NewProjectHandler(projectService)
	projectHandler.LoadRouter(api)
}

func initFile(api *gin.RouterGroup) {
	fileDAO := dao.NewFileDAO(config.DB)
	fileService := service.NewFileService(fileDAO)
	fileHandler := v1.NewFileHandler(fileService)
	fileHandler.LoadRouter(api)
}

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "user.html", gin.H{})
	})

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	initStatic(router)

	api := router.Group("/api")

	initUser(api)
	initProject(api)
	initFile(api)

	return router
}
