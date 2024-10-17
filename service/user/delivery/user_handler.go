package delivery

import (
	"net/http"
	"th3y3m/e-commerce-microservices/service/user/dependency_injection"
	"th3y3m/e-commerce-microservices/service/user/model"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUser(c *gin.Context) {
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
	if err := c.BindJSON(&req); err != nil {
		logrus.Error("Failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "details": err.Error()})
		return
	}

	// Inject the UserUsecase (through dependency injection)
	module := dependency_injection.NewUserUsecaseProvider()

	// Call CreateUser usecase function
	user, err := module.CreateUser(c, &req)
	if err != nil {
		logrus.Error("Failed to create user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "details": err.Error()})
		return
	}

	// Respond with the created user data
	c.JSON(http.StatusOK, user)
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
	if req.Paging.PageIndex == 0 {
		req.Paging.PageIndex = 1
	}
	if req.Paging.PageSize == 0 {
		req.Paging.PageSize = 10
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

func VerifyToken(c *gin.Context) {
	module := dependency_injection.NewUserUsecaseProvider()

	token := c.Query("token")
	userID := c.Query("user_id")

	if token == "" || userID == "" {
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
	isValid, err := module.VerifyToken(c, token, userIDInt)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"is_valid": isValid,
	})
}
