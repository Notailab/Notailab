package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/config"
	"github.com/Notailab/Notailab/internal/dao"
	v1 "github.com/Notailab/Notailab/internal/router/api/v1"
	"github.com/Notailab/Notailab/internal/service"
)

func initCORS(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
}

func initUser(api *gin.RouterGroup) {
	userDAO := dao.NewUserDAO(config.DB)
	userService := service.NewUserService(userDAO)
	userHandler := v1.NewUserHandler(userService)
	userHandler.LoadRouter(api)
}

func initUserSettings(api *gin.RouterGroup) {
	settingDAO := dao.NewUserSettingDAO(config.DB)
	settingService := service.NewUserSettingService(settingDAO)
	settingHandler := v1.NewUserSettingHandler(settingService)
	settingHandler.LoadRouter(api)
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

func initAgent(api *gin.RouterGroup) {
	agentDAO := dao.NewAgentDAO(config.DB)
	projectDAO := dao.NewProjectDAO(config.DB)
	fileDAO := dao.NewFileDAO(config.DB)
	settingDAO := dao.NewUserSettingDAO(config.DB)
	agentService := service.NewAgentService(agentDAO, projectDAO, fileDAO, settingDAO)
	agentHandler := v1.NewAgentHandler(agentService)
	agentHandler.LoadRouter(api)
}

func initStats(api *gin.RouterGroup) {
	statsDAO := dao.NewStatsDAO(config.DB)
	statsService := service.NewStatsService(statsDAO)
	statsHandler := v1.NewStatsHandler(statsService)
	statsHandler.LoadRouter(api)
}

func InitRouter() *gin.Engine {
	router := gin.Default()
	initCORS(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "not found",
		})
	})

	api := router.Group("/api")

	initUser(api)
	initUserSettings(api)
	initProject(api)
	initFile(api)
	initAgent(api)
	initStats(api)

	return router
}
