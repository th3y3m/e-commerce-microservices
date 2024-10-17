package delivery

import (
	"net/http"
	"th3y3m/e-commerce-microservices/service/authentication/dependency_injection"
	"th3y3m/e-commerce-microservices/service/authentication/model"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user model.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module := dependency_injection.NewAuthUsecaseProvider()

	token, err := module.Login(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

func Register(c *gin.Context) {
	var user model.RegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module := dependency_injection.NewAuthUsecaseProvider()

	err := module.RegisterCustomer(user.Email, user.Password, user.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verify your email to complete registration"})
}

func VerifyUserEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	module := dependency_injection.NewAuthUsecaseProvider()

	err := module.VerifyUserEmail(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
