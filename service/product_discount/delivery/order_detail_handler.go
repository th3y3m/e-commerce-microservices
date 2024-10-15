package delivery

import (
	"th3y3m/e-commerce-microservices/service/product_discount/dependency_injection"
	"th3y3m/e-commerce-microservices/service/product_discount/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateProductDiscount(c *gin.Context) {
	var req model.CreateProductDiscountRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewProductDiscountUsecaseProvider()

	productDiscount, err := module.CreateProductDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, productDiscount)
}

func DeleteProductDiscount(c *gin.Context) {
	module := dependency_injection.NewProductDiscountUsecaseProvider()

	var req model.DeleteProductDiscountRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteProductDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ProductDiscount deleted successfully",
	})
}

func GetPaginatedProductDiscount(c *gin.Context) {
	module := dependency_injection.NewProductDiscountUsecaseProvider()

	var req model.GetProductDiscountsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	productDiscounts, err := module.GetProductDiscountList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, productDiscounts)
}
