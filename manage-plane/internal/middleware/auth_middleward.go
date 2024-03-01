package middleware

import (
	"manage-plane/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleward() gin.HandlerFunc {
	// 从请求头中获取token
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": "401",
				"data": "未登录",
				"msg":  "未登录",
			})
			c.Abort()
			return
		}
		token = strings.ReplaceAll(token, "Bearer ", ``)
		clamis, err := utils.ParToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": "401",
				"data": "未登录",
				"msg":  "未登录",
			})
			c.Abort()
			return
		}
		c.Set("userId", clamis["userId"])
		c.Set("role", clamis["sub"])
		c.Set("username", clamis["username"])
		c.Next()
	}
}
