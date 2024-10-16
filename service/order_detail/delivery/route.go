package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	orderDetail := r.Group("/api/orderDetails")
	{
		orderDetail.GET("/GetOrderDetailByID", GetOrderDetailByID)
		orderDetail.GET("", GetOrderDetails)
		orderDetail.POST("", CreateOrderDetail)
		orderDetail.PUT("", UpdateOrderDetail)
		orderDetail.DELETE("", DeleteOrderDetail)
	}

	return r
}
