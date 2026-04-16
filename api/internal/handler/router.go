package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, authHandler AuthHandler) {
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/check", authHandler.CheckUsername)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/signup", authHandler.CreateUser)
			auth.POST("/signup/oauth", authHandler.CreateOAuthUser)
			auth.POST("/login", authHandler.Login)
			auth.POST("/login/oauth", authHandler.OauthLogin)
		}
	}
}
