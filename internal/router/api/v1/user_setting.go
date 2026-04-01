package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/model"
	"github.com/Notailab/Notailab/internal/service"
	"github.com/Notailab/Notailab/middleware/auth"
)

type UserSettingHandler struct {
	settingService *service.UserSettingService
}

func NewUserSettingHandler(settingService *service.UserSettingService) *UserSettingHandler {
	return &UserSettingHandler{settingService: settingService}
}

type UserSettingUpdateRequest struct {
	LLMProvider string  `json:"llm_provider"`
	LLMBaseURL  string  `json:"llm_base_url"`
	LLMAPIKey   string  `json:"llm_api_key"`
	LLMModel    string  `json:"llm_model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   *int    `json:"max_tokens"`
}

func (h *UserSettingHandler) GetSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未获取到用户信息，请重新登录"})
		return
	}

	setting, err := h.settingService.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": fmt.Sprintf("获取用户设置失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    setting,
	})
}

func (h *UserSettingHandler) UpdateSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未获取到用户信息，请重新登录"})
		return
	}

	var req UserSettingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	setting := &model.UserSetting{
		LLMProvider: req.LLMProvider,
		LLMBaseURL:  req.LLMBaseURL,
		LLMAPIKey:   req.LLMAPIKey,
		LLMModel:    req.LLMModel,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	if err := h.settingService.Upsert(userID.(uint), setting); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": fmt.Sprintf("保存用户设置失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data":    setting,
	})
}

func (h *UserSettingHandler) LoadRouter(api *gin.RouterGroup) {
	setting := api.Group("/user", auth.AuthMiddleware())
	{
		setting.GET("/settings", h.GetSettings)
		setting.PUT("/settings", h.UpdateSettings)
	}
}
