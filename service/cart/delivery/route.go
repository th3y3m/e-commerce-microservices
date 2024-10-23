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

		cart.POST("/delete-unit-item", DeleteUnitItem)
		cart.POST("/save-cart-to-cookie-handler", SaveCartToCookieHandler)
		cart.DELETE("/remove-from-cart", RemoveFromCart)
		cart.DELETE("/delete-cart-in-cookie", DeleteCartInCookie)
		cart.GET("/number-of-items-in-cart-cookie", NumberOfItemsInCartCookie)

	}

	return r
}
