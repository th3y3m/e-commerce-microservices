package delivery

import (
	"net/http"
	"strconv"
	"th3y3m/e-commerce-microservices/service/momo/dependency_injection"

	"github.com/gin-gonic/gin"
)

func CreateMoMoUrl(c *gin.Context) {
	module := dependency_injection.NewMoMoUsecaseProvider()

	amountStr := c.Query("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	orderID := c.Query("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	paymentUrl, err := module.CreateMoMoUrl(amount, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_url": paymentUrl})
}

func ValidateMoMoResponse(c *gin.Context) {
	module := dependency_injection.NewMoMoUsecaseProvider()

	queryParams := c.Request.URL.Query()
	res, err := module.ValidateMoMoResponse(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res})
}
