package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	vnpay := r.Group("/api/vnpay")
	{
		vnpay.POST("", CreateVnPayUrl)
		vnpay.GET("/validate", ValidateVnPayResponse)
	}

	return r
}
