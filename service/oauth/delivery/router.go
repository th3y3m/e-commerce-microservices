package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/google/callback", GoogleCallback)
		auth.POST("/facebook/callback", FacebookCallback)
		auth.GET("/google/login", GoogleLogin)
		auth.GET("/google/login", FacebookLogin)
	}

	return r
}
