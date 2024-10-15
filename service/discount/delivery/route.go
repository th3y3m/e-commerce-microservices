package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	discount := r.Group("/api/discounts")
	{
		discount.GET("/:discount_id", GetDiscountByID)
		discount.POST("/", CreateDiscount)
		discount.PUT("/", UpdateDiscount)
		discount.DELETE("/", DeleteDiscount)
	}

	return r
}
