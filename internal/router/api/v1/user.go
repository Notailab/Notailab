package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/service"
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

	isValid, err := h.userService.UserLogin(username, password)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  404,
			"error": "User not found",
		})
		return
	}

	if !isValid {
		c.JSON(200, gin.H{
			"code":  401,
			"error": "Invalid credentials",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":     200,
		"message":  "User logged in successfully",
		"username": username,
	})
}

func (h *UserHandler) LoadRouter(e *gin.Engine) {
	e.POST("/api/user/register", h.Register)
	e.POST("/api/user/login", h.Login)
}

