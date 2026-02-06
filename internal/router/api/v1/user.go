package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}
	username := req.Username
	password := req.Password

	err := h.userService.CreateUser(username, password)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to register user: %v", err)})
		return
	}
	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}
	username := req.Username
	password := req.Password

	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if user.Password != password {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	}
	c.JSON(200, gin.H{
		"message":  "User logged in successfully",
		"username": user.Username,
		"password": user.Password,
	})
}

func (h *UserHandler) LoadRouter(e *gin.Engine) {
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
}