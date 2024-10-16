package delivery

import (
	"strconv"
	"th3y3m/e-commerce-microservices/service/discount/dependency_injection"
	"th3y3m/e-commerce-microservices/service/discount/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetDiscountByID(c *gin.Context) {
	id := c.Param("discount_id")

	module := dependency_injection.NewDiscountUsecaseProvider()

	var req model.GetDiscountRequest
	discountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid discount ID",
		})
		return
	}
	req.DiscountID = discountID

	discount, err := module.GetDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, discount)
}

func GetAllDiscounts(c *gin.Context) {
	module := dependency_injection.NewDiscountUsecaseProvider()

	discounts, err := module.GetAllDiscounts(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, discounts)
}

func CreateDiscount(c *gin.Context) {
	var req model.CreateDiscountRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewDiscountUsecaseProvider()

	discount, err := module.CreateDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, discount)
}

func UpdateDiscount(c *gin.Context) {
	var req model.UpdateDiscountRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewDiscountUsecaseProvider()

	discount, err := module.UpdateDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, discount)
}

func DeleteDiscount(c *gin.Context) {
	module := dependency_injection.NewDiscountUsecaseProvider()

	var req model.DeleteDiscountRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Discount deleted successfully",
	})
}
