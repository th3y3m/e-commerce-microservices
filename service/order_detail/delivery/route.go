package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	orderDetail := r.Group("/api/orderDetails")
	{
		orderDetail.GET("/:orderDetail_id", GetOrderDetailByID)
		orderDetail.GET("/", GetPaginatedOrderDetail)
		orderDetail.POST("/", CreateOrderDetail)
		orderDetail.PUT("/", UpdateOrderDetail)
		orderDetail.DELETE("/", DeleteOrderDetail)
	}

	return r
}
