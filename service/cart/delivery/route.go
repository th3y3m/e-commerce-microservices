package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	cart := r.Group("/api/carts")
	{
		cart.GET("/:cart_id", GetCartByID)
		cart.POST("", CreateCart)
		cart.PUT("", UpdateCart)
		cart.DELETE("", DeleteCart)
		cart.GET("/get-user-cart/:user_id", GetUserCart)
		cart.POST("/add-item", AddProductToShoppingCart)
		cart.PUT("/remove-item", RemoveProductFromShoppingCart)
	}

	return r
}
