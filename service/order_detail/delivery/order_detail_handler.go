package delivery

import (
	"th3y3m/e-commerce-microservices/service/order_detail/dependency_injection"
	"th3y3m/e-commerce-microservices/service/order_detail/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetOrderDetailByID(c *gin.Context) {
	module := dependency_injection.NewOrderDetailUsecaseProvider()

	var req model.GetOrderDetailRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	orderDetail, err := module.GetOrderDetail(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, orderDetail)
}

func CreateOrderDetail(c *gin.Context) {
	var req model.CreateOrderDetailRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewOrderDetailUsecaseProvider()

	orderDetail, err := module.CreateOrderDetail(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, orderDetail)
}

func UpdateOrderDetail(c *gin.Context) {
	var req model.UpdateOrderDetailRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewOrderDetailUsecaseProvider()

	orderDetail, err := module.UpdateOrderDetail(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, orderDetail)
}

func DeleteOrderDetail(c *gin.Context) {
	module := dependency_injection.NewOrderDetailUsecaseProvider()

	var req model.DeleteOrderDetailRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteOrderDetail(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "OrderDetail deleted successfully",
	})
}

func GetOrderDetails(c *gin.Context) {
	module := dependency_injection.NewOrderDetailUsecaseProvider()

	var req model.GetOrderDetailsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	orderDetails, err := module.GetOrderDetailList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, orderDetails)
}
