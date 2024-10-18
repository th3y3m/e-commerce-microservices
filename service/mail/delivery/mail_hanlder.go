package delivery

import (
	"strconv"
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

// func SendOrderDetails(c *gin.Context) {
// 	var request model.SendOrderDetailsRequest

// 	// Bind the incoming JSON body to the request struct
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
// 		return
// 	}
// 	module := dependency_injection.NewMailUsecaseProvider()
// 	// Call the use case to send the order details email
// 	err := module.SendOrderDetails(request.Customer, request.Order, request.OrderDetails)
// 	if err != nil {
// 		// Log and return an error if the email sending failed
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send order details", "details": err.Error()})
// 		return
// 	}

// 	// Respond with success
// 	c.JSON(http.StatusOK, gin.H{"message": "Order details sent successfully"})
// }

func SendNotification(c *gin.Context) {
	orderIDStr := c.Query("order_id")
	url := c.Query("url")

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid order ID",
		})
		return
	}

	module := dependency_injection.NewMailUsecaseProvider()

	err = module.SendNotification(c, orderID, url)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to send notification",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Notification sent",
	})
}
