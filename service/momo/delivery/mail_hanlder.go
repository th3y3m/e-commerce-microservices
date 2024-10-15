package delivery

import (
	"th3y3m/e-commerce-microservices/service/mail/dependency_injection"

	"github.com/gin-gonic/gin"
)

func SendMail(c *gin.Context) {
	to := c.Query("to")
	token := c.Query("token")

	if to == "" || token == "" {
		c.JSON(400, gin.H{
			"message": "Missing required fields",
		})
		return
	}

	module := dependency_injection.NewMailUsecaseProvider()

	err := module.SendMail(to, token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to send mail",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Mail sent",
	})
}
