package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	cartItem := r.Group("/api/cartItems")
	{
		cartItem.GET("/:cart_id/:product_id", GetCartItemByID)
		cartItem.GET("/", GetCartItems)
		cartItem.POST("/", CreateCartItem)
		cartItem.PUT("/", UpdateCartItem)
		cartItem.DELETE("/", DeleteCartItem)
	}

	return r
}
