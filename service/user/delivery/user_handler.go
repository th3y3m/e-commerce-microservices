package delivery

import (
	"th3y3m/e-commerce-microservices/service/user/dependency_injection"
	"th3y3m/e-commerce-microservices/service/user/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUserByID(c *gin.Context) {
	module := dependency_injection.NewUserUsecaseProvider()

	var req model.GetUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	user, err := module.GetUser(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, user)
}

func GetAllUsers(c *gin.Context) {
	module := dependency_injection.NewUserUsecaseProvider()

	users, err := module.GetAllUsers(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, users)
}

func CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewUserUsecaseProvider()

	user, err := module.CreateUser(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	var req model.UpdateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewUserUsecaseProvider()

	user, err := module.UpdateUser(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	module := dependency_injection.NewUserUsecaseProvider()

	var req model.DeleteUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteUser(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User deleted successfully",
	})
}

func GetPaginatedUser(c *gin.Context) {
	module := dependency_injection.NewUserUsecaseProvider()

	var req model.GetUsersRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	users, err := module.GetUserList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, users)
}
