package delivery

import (
	"strconv"
	"th3y3m/e-commerce-microservices/service/cart/dependency_injection"
	"th3y3m/e-commerce-microservices/service/cart/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCartByID(c *gin.Context) {
	var req model.GetCartRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	cart, err := module.GetCart(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cart)
}

func CreateCart(c *gin.Context) {
	var req model.CreateCartRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	cart, err := module.CreateCart(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cart)
}

func UpdateCart(c *gin.Context) {
	var req model.UpdateCartRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	cart, err := module.UpdateCart(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cart)
}

func DeleteCart(c *gin.Context) {
	var req model.DeleteCartRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	err = module.DeleteCart(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Cart deleted successfully",
	})
}

func GetUserCart(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	cart, err := module.GetUserCart(c, userID)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cart)
}

func AddProductToShoppingCart(c *gin.Context) {
	var req model.AddProductToShoppingCartRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	err = module.AddProductToShoppingCart(c, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{"message": "Product added to shopping cart successfully"})
}

func RemoveProductFromShoppingCart(c *gin.Context) {
	var req model.RemoveProductFromShoppingCartRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartUsecaseProvider()

	err = module.RemoveProductFromShoppingCart(c, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{"message": "Product removed from shopping cart successfully"})
}

func DeleteUnitItem(c *gin.Context) {
	// Initialize cart usecase module
	module := dependency_injection.NewCartUsecaseProvider()

	// Extract productId from URL or request
	productIDParam := c.Param("productId") // Assuming productId is part of the URL
	productId, err := strconv.ParseInt(productIDParam, 10, 64)
	if err != nil {
		logrus.Error("Invalid product ID")
		c.JSON(400, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Call the DeleteUnitItem method, passing the required Gin context objects
	err = module.DeleteUnitItem(c.Writer, c.Request, productId)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Return success response
	c.JSON(200, gin.H{"message": "Unit item deleted successfully"})
}

func RemoveFromCart(c *gin.Context) {
	// Initialize cart usecase module
	module := dependency_injection.NewCartUsecaseProvider()

	// Extract productId from URL or request
	productIDParam := c.Query("productId") // Assuming productId is part of the URL
	productId, err := strconv.ParseInt(productIDParam, 10, 64)
	if err != nil {
		logrus.Error("Invalid product ID")
		c.JSON(400, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Call the RemoveFromCart method, passing the required Gin context objects
	err = module.RemoveFromCart(c.Writer, c.Request, productId)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Return success response
	c.JSON(200, gin.H{"message": "Product removed from cart successfully"})
}

func DeleteCartInCookie(c *gin.Context) {
	// Initialize cart usecase module
	module := dependency_injection.NewCartUsecaseProvider()

	// Call the DeleteCartInCookie method, passing the required Gin context objects
	err := module.DeleteCartInCookie(c.Writer)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Return success response
	c.JSON(200, gin.H{"message": "Cart deleted from cookie successfully"})
}

func NumberOfItemsInCartCookie(c *gin.Context) {
	// Initialize cart usecase module
	module := dependency_injection.NewCartUsecaseProvider()

	// Call the NumberOfItemsInCartCookie method, passing the required Gin context objects
	numItems, err := module.NumberOfItemsInCartCookie(c.Request)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Return the number of items in the cart as a JSON response
	c.JSON(200, gin.H{"numItems": numItems})
}

func SaveCartToCookieHandler(c *gin.Context) {
	// Initialize cart usecase module
	module := dependency_injection.NewCartUsecaseProvider()

	// Extract productId from URL or request
	productIDParam := c.Query("productId") // Assuming productId is part of the URL
	productId, err := strconv.ParseInt(productIDParam, 10, 64)
	if err != nil {
		logrus.Error("Invalid product ID")
		c.JSON(400, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Call the SaveCartToCookieHandler method, passing the required Gin context objects
	err = module.SaveCartToCookieHandler(c.Writer, c.Request, productId)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Return success response
	c.JSON(200, gin.H{"message": "Cart saved to cookie successfully"})
}
