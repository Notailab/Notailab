package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/Notailab/Notailab/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "未登录，请先登录",
			})
			c.Abort()
			return
		}

		// 去掉 Bearer 前缀（如果有）
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		// 2. 解析 Token
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 3. 将用户信息存入上下文，供后续接口使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
