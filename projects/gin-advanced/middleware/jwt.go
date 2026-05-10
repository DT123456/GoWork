package middleware

import (
	"gin-advanced/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		// 支持两种 token 传递方式
		// 1. Authorization: Bearer xxx (标准格式，前端使用)
		// 2. token: xxx (备用)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 从 "Bearer xxx" 中提取 token
			if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
				token = authHeader[7:]
			}
		}
		// 如果 Authorization 头没有，尝试 token 头
		if token == "" {
			token = c.GetHeader("token")
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登录"})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "登录已过期"})
			c.Abort()
			return
		}
		// 将用户 ID 存入 context 供后续使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
