package middleware

import (
	"planet/internal/pkg"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := pkg.ParseAccessToken(token)
		if err != nil {
			pkg.Fail(c, 401, "인증이 필요합니다")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			claims, err := pkg.ParseAccessToken(token)
			if err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
			}
		}
		c.Next()
	}
}
