package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	product := r.Group("/api/products")
	{
		product.GET("/:product_id", GetProductByID)
		product.GET("/", GetPaginatedProduct)
		product.POST("/", CreateProduct)
		product.PUT("/:product_id", UpdateProduct)
		product.DELETE("/:product_id", DeleteProduct)
	}

	return r
}
