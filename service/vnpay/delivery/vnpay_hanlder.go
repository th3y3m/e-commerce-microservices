package delivery

import (
	"net/http"
	"strconv"
	"th3y3m/e-commerce-microservices/service/vnpay/dependency_injection"

	"github.com/gin-gonic/gin"
)

func CreateVnPayUrl(c *gin.Context) {
	VnPayConfig := dependency_injection.NewVnpayUsecaseProvider()

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

	paymentUrl, err := VnPayConfig.CreateVNPayUrl(amount, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_url": paymentUrl})
}

func ValidateVnPayResponse(c *gin.Context) {
	VnPayConfig := dependency_injection.NewVnpayUsecaseProvider()

	queryParams := c.Request.URL.Query()
	res, err := VnPayConfig.ValidateVNPayResponse(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res})
}
