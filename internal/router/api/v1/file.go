package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/model"
	"github.com/Notailab/Notailab/internal/service"
	"github.com/Notailab/Notailab/middleware/auth"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

type FileCreateRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Content   string `json:"content" binding:"required"`
	IsHidden  bool   `json:"is_hidden"`
}

func (h *FileHandler) CreateFile(c *gin.Context) {
	var req FileCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	err := h.fileService.CreateFile(&model.File{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Content:   req.Content,
		IsHidden:  req.IsHidden || strings.HasPrefix(req.Name, "."),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "Failed to create file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File created successfully",
	})
}

func (h *FileHandler) GetFilesByProjectID(c *gin.Context) {
	// 换成 POST 请求，从 JSON body 中获取 project_id
	var req struct {
		Project_id uint `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	files, err := h.fileService.GetFilesByProjectID(req.Project_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "Failed to get files: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Files retrieved successfully",
		"data":    files,
	})
}

func (h *FileHandler) UpdateFileContent(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Missing file_id path parameter",
		})
		return
	}

	var fileID uint
	if _, err := fmt.Sscanf(fileIDStr, "%d", &fileID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid file_id path parameter",
		})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	err := h.fileService.UpdateFileContent(fileID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "Failed to update file content: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File content updated successfully",
	})
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Missing file_id path parameter",
		})
		return
	}

	var fileID uint
	if _, err := fmt.Sscanf(fileIDStr, "%d", &fileID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid file_id path parameter",
		})
		return
	}

	err := h.fileService.DeleteFile(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "Failed to delete file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File deleted successfully",
	})
}

func (h *FileHandler) LoadRouter(rg *gin.RouterGroup) {
	file := rg.Group("/file", auth.AuthMiddleware())
	{
		file.POST("/new", h.CreateFile)
		file.POST("/all", h.GetFilesByProjectID)
		file.PUT("/:file_id", h.UpdateFileContent)
		file.DELETE("/:file_id", h.DeleteFile)
	}
}
