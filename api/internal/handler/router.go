package handler

import (
	"planet/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	authHandler AuthHandler,
	taskHandler TaskHandler,
	userHandler UserHandler,
) {

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
		tasks.Use(middleware.AuthMiddleware())
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.POST("/:id/toggle", taskHandler.ToggleTask)
		}

		users := v1.Group("/users")
		users.Use(middleware.OptionalAuthMiddleware())
		{
			users.GET("/:userid", userHandler.GetProfile)
			users.GET("/:userid/tasks", taskHandler.GetTasksByMonth)
		}
	}
}
