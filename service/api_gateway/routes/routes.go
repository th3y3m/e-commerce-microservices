package routes

import (
	"th3y3m/e-commerce-microservices/service/api_gateway/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the routes for the API Gateway
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Product Service Routes (handled by API Gateway)
	product := router.Group("/api/products")
	{
		product.GET("/:product_id", handler.GetProductByID)
		product.GET("/", handler.GetPaginatedProducts)
		product.POST("/", handler.CreateProduct)
		product.PUT("/:product_id", handler.UpdateProduct)
		product.DELETE("/:product_id", handler.DeleteProduct)
	}

	// Add routes for other services (e.g., cart, payment, etc.) here

	return router
}
