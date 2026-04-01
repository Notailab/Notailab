package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/service"
	"github.com/Notailab/Notailab/middleware/auth"
)

type AgentHandler struct {
	agentService *service.AgentService
}

func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

type AgentChatResponse struct {
	ConversationID uint                `json:"conversation_id"`
	Reply          string              `json:"reply"`
	History        []map[string]string `json:"history"`
}

func (h *AgentHandler) Chat(c *gin.Context) {
	var req service.AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息，请重新登录",
		})
		return
	}

	result, err := h.agentService.Chat(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("Agent chat failed: %v", err),
		})
		return
	}

	history := make([]map[string]string, 0, len(result.History))
	for _, item := range result.History {
		history = append(history, map[string]string{
			"role":    item.Role,
			"content": item.Content,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": AgentChatResponse{
			ConversationID: result.Conversation.ConversationID,
			Reply:          result.Reply,
			History:        history,
		},
	})
}

func (h *AgentHandler) LoadRouter(api *gin.RouterGroup) {
	agent := api.Group("/agent", auth.AuthMiddleware())
	{
		agent.POST("/chat", h.Chat)
	}
}
