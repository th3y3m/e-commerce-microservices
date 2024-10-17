package delivery

import (
	"context"
	"log"
	"th3y3m/e-commerce-microservices/service/product/dependency_injection"
	"th3y3m/e-commerce-microservices/service/product/rabbitmq"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()
	module := dependency_injection.NewProductUsecaseProvider()
	ctx := context.Background()

	go func() {
		if err := rabbitmq.ConsumeInventoryUpdates(ctx, module); err != nil {
			log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	product := r.Group("/api/products")
	{
		product.GET("/:product_id", GetProductByID)
		product.GET("", GetPaginatedProduct)
		product.POST("", CreateProduct)
		product.PUT("", UpdateProduct)
		product.DELETE("", DeleteProduct)
	}

	return r
}
