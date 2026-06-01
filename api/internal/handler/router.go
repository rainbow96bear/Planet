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
	searchHandler SearchHandler,
	feedHandler FeedHandler,
	notificationHandler NotificationHandler,
	reactionHandler ReactionHandler,
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
			tasks.DELETE("/:task_id", taskHandler.DeleteTask)
			tasks.POST("/:task_id/toggle", taskHandler.ToggleTask)

			tasks.POST("/:task_id/reactions", reactionHandler.AddReaction)
			tasks.DELETE("/:task_id/reactions", reactionHandler.RemoveReaction)
		}

		users := v1.Group("/users")
		users.Use(middleware.OptionalAuthMiddleware())
		{
			users.GET("/:userid", userHandler.GetProfile)
			users.GET("/:userid/tasks", taskHandler.GetTasksByMonth)
		}

		usersAuth := v1.Group("/users")
		usersAuth.Use(middleware.AuthMiddleware())
		{
			usersAuth.PATCH("/:userid", userHandler.UpdateProfile)
			usersAuth.POST("/:userid/follow", userHandler.Follow)
			usersAuth.DELETE("/:userid/follow", userHandler.Unfollow)
		}

		search := v1.Group("/search")
		search.Use(middleware.OptionalAuthMiddleware())
		{
			search.GET("/users", searchHandler.SearchUsers)
		}

		feed := v1.Group("/feed")
		feed.Use(middleware.OptionalAuthMiddleware())
		{
			feed.GET("", feedHandler.GetFeed)
			feed.GET("/explore", feedHandler.GetExploreFeed)
		}

		notifications := v1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware())
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PATCH("/read-all", notificationHandler.MarkAllAsRead)
		}
	}
}
