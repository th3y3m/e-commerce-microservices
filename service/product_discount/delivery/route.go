package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	productDiscount := r.Group("/api/productDiscounts")
	{
		productDiscount.GET("/", GetPaginatedProductDiscount)
		productDiscount.POST("/", CreateProductDiscount)
		productDiscount.DELETE("/", DeleteProductDiscount)
	}

	return r
}
