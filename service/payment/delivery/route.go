package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	payment := r.Group("/api/payments")
	{
		payment.GET("/:payment_id", GetPaymentByID)
		payment.GET("", GetPaginatedPayment)
		payment.POST("", CreatePayment)
		payment.PUT("", UpdatePayment)
	}

	return r
}
