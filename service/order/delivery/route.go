package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	order := r.Group("/api/orders")
	{
		order.GET("/:order_id", GetOrderByID)
		order.GET("/", GetPaginatedOrder)
		order.POST("/", CreateOrder)
		order.PUT("/", UpdateOrder)
		order.DELETE("/", DeleteOrder)
	}

	return r
}
