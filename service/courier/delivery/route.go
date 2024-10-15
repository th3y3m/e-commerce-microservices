package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	courier := r.Group("/api/couriers")
	{
		courier.GET("/:courier_id", GetCourierByID)
		courier.POST("/", CreateCourier)
		courier.PUT("/", UpdateCourier)
		courier.DELETE("/", DeleteCourier)
	}

	return r
}
