package delivery

import (
	"th3y3m/e-commerce-microservices/service/courier/dependency_injection"
	"th3y3m/e-commerce-microservices/service/courier/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCourierByID(c *gin.Context) {
	var req model.GetCourierRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCourierUsecaseProvider()

	courier, err := module.GetCourier(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, courier)
}

func GetAllCouriers(c *gin.Context) {
	module := dependency_injection.NewCourierUsecaseProvider()

	couriers, err := module.GetAllCouriers(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, couriers)
}

func CreateCourier(c *gin.Context) {
	var req model.CreateCourierRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCourierUsecaseProvider()

	courier, err := module.CreateCourier(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, courier)
}

func UpdateCourier(c *gin.Context) {
	var req model.UpdateCourierRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCourierUsecaseProvider()

	courier, err := module.UpdateCourier(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, courier)
}

func DeleteCourier(c *gin.Context) {
	module := dependency_injection.NewCourierUsecaseProvider()

	var req model.DeleteCourierRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteCourier(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Courier deleted successfully",
	})
}
