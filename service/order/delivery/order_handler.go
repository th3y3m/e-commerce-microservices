package delivery

import (
	"th3y3m/e-commerce-microservices/service/order/dependency_injection"
	"th3y3m/e-commerce-microservices/service/order/model"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetOrderByID(c *gin.Context) {
	id := c.Param("order_id")

	module := dependency_injection.NewOrderUsecaseProvider()

	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid Order ID",
		})
		return
	}

	var req model.GetOrderRequest
	req.OrderID = orderID

	order, err := module.GetOrder(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, order)
}

func GetAllOrders(c *gin.Context) {
	module := dependency_injection.NewOrderUsecaseProvider()

	orders, err := module.GetAllOrders(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, orders)
}

func CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewOrderUsecaseProvider()

	order, err := module.CreateOrder(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, order)
}

func UpdateOrder(c *gin.Context) {
	var req model.UpdateOrderRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewOrderUsecaseProvider()

	order, err := module.UpdateOrder(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, order)
}

func DeleteOrder(c *gin.Context) {
	module := dependency_injection.NewOrderUsecaseProvider()

	var req model.DeleteOrderRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteOrder(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Order deleted successfully",
	})
}

func GetPaginatedOrder(c *gin.Context) {
	module := dependency_injection.NewOrderUsecaseProvider()

	var req model.GetOrdersRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	if req.Paging.PageIndex == 0 {
		req.Paging.PageIndex = 1
	}
	if req.Paging.PageSize == 0 {
		req.Paging.PageSize = 10
	}

	orders, err := module.GetOrderList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, orders)
}
