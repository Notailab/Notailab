package router

import (
	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/config"
	"github.com/Notailab/Notailab/internal/dao"
	v1 "github.com/Notailab/Notailab/internal/router/api/v1"
	"github.com/Notailab/Notailab/internal/service"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	
	userDAO := dao.NewUserDAO(config.DB)
	userService := service.NewUserService(userDAO)
	userHandler := v1.NewUserHandler(userService)
	userHandler.LoadRouter(router)

	return router
}