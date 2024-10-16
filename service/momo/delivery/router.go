package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	momo := r.Group("/api/momo")
	{
		momo.POST("/", CreateMoMoUrl)
		momo.GET("/validate", ValidateMoMoResponse)
	}

	return r
}
