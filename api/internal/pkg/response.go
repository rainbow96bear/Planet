package pkg

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"data":  data,
		"error": nil,
	})
}

func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"data":  nil,
		"error": gin.H{"message": message},
	})
}
