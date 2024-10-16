package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	cartItem := r.Group("/api/cartItems")
	{
		cartItem.GET("/GetCartItemByID", GetCartItemByID)
		cartItem.GET("/", GetCartItems)
		cartItem.POST("/", CreateCartItem)
		cartItem.PUT("/", UpdateCartItem)
		cartItem.PUT("/UpdateOrCreateCartItem", UpdateOrCreateCartItem)
		cartItem.DELETE("/", DeleteCartItem)
	}

	return r
}
