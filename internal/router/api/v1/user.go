package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/service"
	"github.com/Notailab/Notailab/middleware/auth"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfoRequest struct {
	UserId   uint   `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":  400,
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	username := req.Username
	password := req.Password
	email := req.Email

	if !h.userService.CheckUsername(username) {
		c.JSON(http.StatusOK, gin.H{
			"code":    1001,
			"message": "Username is already exist.",
		})
		return
	}

	if !h.userService.CheckEmail(email) {
		c.JSON(http.StatusOK, gin.H{
			"code":    1002,
			"message": "Email is already exist.",
		})
		return
	}

	err := h.userService.CreateUser(username, password, email)
	if err != nil {
		c.JSON(500, gin.H{
			"code":  500,
			"error": fmt.Sprintf("Failed to register user: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":  400,
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	username := req.Username
	password := req.Password

	token, err := h.userService.UserLogin(username, password)
	if err != nil || token == "" {
		c.JSON(200, gin.H{
			"code":  1000,
			"error": "Login failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":     200,
		"message":  "User logged in successfully",
		"username": username,
		"token":    token,
	})
}

func (h *UserHandler) Info(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		// 未取到 user_id（理论上不会走到这里，因为中间件已校验 Token 并存入）
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息，请重新登录",
		})
		return
	}

	user, err := h.userService.GetUserByID(user_id.(uint))
	if err != nil {
		c.JSON(200, gin.H{
			"code":  404,
			"error": "No Authorzation",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":     200,
		"id":       user_id,
		"message":  "Authorzation successfully",
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
	})
}

func (h *UserHandler) LoadRouter(api *gin.RouterGroup) {
	api.POST("/user/register", h.Register)
	api.POST("/user/login", h.Login)
	user := api.Group("/user", auth.AuthMiddleware())
	{
		user.POST("/info", h.Info)
	}
}
