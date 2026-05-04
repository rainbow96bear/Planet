package middleware

import (
	"planet/internal/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := pkg.ParseAccessToken(tokenString)
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
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := pkg.ParseAccessToken(tokenString)
			if err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
			}
		}
		c.Next()
	}
}
