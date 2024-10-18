package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.GET("/google/callback", GoogleCallback)
		auth.GET("/facebook/callback", FacebookCallback)
		auth.GET("/google/login", GoogleLogin)
		auth.GET("/facebook/login", FacebookLogin)
	}

	return r
}
