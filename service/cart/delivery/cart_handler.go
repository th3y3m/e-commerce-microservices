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
