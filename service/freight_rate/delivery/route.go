package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	freightRate := r.Group("/api/freightRates")
	{
		freightRate.GET("/:freightRate_id", GetFreightRateByID)
		freightRate.POST("", CreateFreightRate)
		freightRate.PUT("", UpdateFreightRate)
		freightRate.DELETE("", DeleteFreightRate)
	}

	return r
}
