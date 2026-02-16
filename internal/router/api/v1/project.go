package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/model"
	"github.com/Notailab/Notailab/internal/service"
	"github.com/Notailab/Notailab/middleware/auth"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

type ProjectCreateRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req ProjectCreateRequest
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息，请重新登录",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":  400,
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	if !h.projectService.CheckTitleByUserId(userId.(uint), req.Title) {
		c.JSON(http.StatusOK, gin.H{
			"code":    1001,
			"message": "Title is already exist.",
		})
		return
	}

	err := h.projectService.CreateProject(&model.Project{
		UserID:      userId.(uint),
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		c.JSON(500, gin.H{
			"code":  500,
			"error": fmt.Sprintf("Failed to create project: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Create project successfully",
	})
}

func (h *ProjectHandler) GetProjectTitles(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息，请重新登录",
		})
		return
	}

	titles, err := h.projectService.GetAllProjectTitlesByUserId(userId.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success!",
		"titles":  titles,
	})
}

func (h *ProjectHandler) LoadRouter(api *gin.RouterGroup) {
	user := api.Group("/project", auth.AuthMiddleware())
	{
		user.POST("/new", h.CreateProject)
		user.POST("/titles", h.GetProjectTitles)
	}
}
