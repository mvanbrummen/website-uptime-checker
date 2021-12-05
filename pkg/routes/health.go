package routes

import "github.com/gin-gonic/gin"

func HealthRoute(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
