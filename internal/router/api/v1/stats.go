package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/internal/service"
	"github.com/Notailab/Notailab/middleware/auth"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) Overview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息，请重新登录",
		})
		return
	}

	overview, err := h.statsService.GetOverview(userID.(uint))
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
		"data":    overview,
	})
}

func (h *StatsHandler) LoadRouter(api *gin.RouterGroup) {
	stats := api.Group("/stats", auth.AuthMiddleware())
	{
		stats.GET("/overview", h.Overview)
	}
}
