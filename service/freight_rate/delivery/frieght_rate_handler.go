package delivery

import (
	"th3y3m/e-commerce-microservices/service/freight_rate/dependency_injection"
	"th3y3m/e-commerce-microservices/service/freight_rate/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetFreightRateByID(c *gin.Context) {
	module := dependency_injection.NewFreightRateUsecaseProvider()

	var req model.GetFreightRateRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
	freightRate, err := module.GetFreightRate(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, freightRate)
}

func GetAllFreightRates(c *gin.Context) {
	module := dependency_injection.NewFreightRateUsecaseProvider()

	freightRates, err := module.GetAllFreightRates(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, freightRates)
}

func CreateFreightRate(c *gin.Context) {
	var req model.CreateFreightRateRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewFreightRateUsecaseProvider()

	freightRate, err := module.CreateFreightRate(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, freightRate)
}

func UpdateFreightRate(c *gin.Context) {
	var req model.UpdateFreightRateRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewFreightRateUsecaseProvider()

	freightRate, err := module.UpdateFreightRate(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, freightRate)
}

func DeleteFreightRate(c *gin.Context) {
	module := dependency_injection.NewFreightRateUsecaseProvider()

	var req model.DeleteFreightRateRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteFreightRate(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "FreightRate deleted successfully",
	})
}
