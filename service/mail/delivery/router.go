package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	mail := r.Group("/api/mail")
	{
		mail.POST("/send-mail", SendMail)
		mail.POST("/send-order-details", SendOrderDetails)
	}

	return r
}
