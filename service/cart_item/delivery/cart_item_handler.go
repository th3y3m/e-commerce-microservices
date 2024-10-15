package delivery

import (
	"th3y3m/e-commerce-microservices/service/cart_item/dependency_injection"
	"th3y3m/e-commerce-microservices/service/cart_item/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCartItemByID(c *gin.Context) {
	var req model.GetCartItemRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartItemUsecaseProvider()

	cartItem, err := module.GetCartItem(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cartItem)
}

func CreateCartItem(c *gin.Context) {
	var req model.CreateCartItemRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartItemUsecaseProvider()

	cartItem, err := module.CreateCartItem(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cartItem)
}

func UpdateCartItem(c *gin.Context) {
	var req model.UpdateCartItemRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartItemUsecaseProvider()

	cartItem, err := module.UpdateCartItem(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cartItem)
}

func DeleteCartItem(c *gin.Context) {
	var req model.DeleteCartItemRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartItemUsecaseProvider()

	err = module.DeleteCartItem(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "CartItem deleted successfully",
	})
}

func GetCartItems(c *gin.Context) {
	var req model.GetCartItemsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCartItemUsecaseProvider()

	cartItems, err := module.GetCartItemList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, cartItems)
}
