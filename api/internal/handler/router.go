package handler

import (
	"planet/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authHandler AuthHandler, taskHandler TaskHandler) {
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

		tasks := v1.Group("/tasks")
		{
			// 인증 필요 (본인만)
			authenticated := tasks.Group("")
			authenticated.Use(middleware.AuthMiddleware())
			{
				authenticated.POST("", taskHandler.CreateTask)
				authenticated.DELETE("/:id", taskHandler.DeleteTask)
			}

		}
	}
}
