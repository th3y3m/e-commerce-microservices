package delivery

import (
	"th3y3m/e-commerce-microservices/service/payment/dependency_injection"
	"th3y3m/e-commerce-microservices/service/payment/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetPaymentByID(c *gin.Context) {
	module := dependency_injection.NewPaymentUsecaseProvider()

	var req model.GetPaymentRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	payment, err := module.GetPayment(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, payment)
}

func GetAllPayments(c *gin.Context) {
	module := dependency_injection.NewPaymentUsecaseProvider()

	payments, err := module.GetAllPayments(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, payments)
}

func CreatePayment(c *gin.Context) {
	var req model.CreatePaymentRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewPaymentUsecaseProvider()

	payment, err := module.CreatePayment(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, payment)
}

func UpdatePayment(c *gin.Context) {
	var req model.UpdatePaymentRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewPaymentUsecaseProvider()

	payment, err := module.UpdatePayment(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, payment)
}

func GetPaginatedPayment(c *gin.Context) {
	module := dependency_injection.NewPaymentUsecaseProvider()

	var req model.GetPaymentsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	payments, err := module.GetPaymentList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, payments)
}
