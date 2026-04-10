package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, userHandler UserHandler) {
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
		}
	}
}
