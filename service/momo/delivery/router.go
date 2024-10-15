package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	mail := r.Group("/api/mail")
	{
		mail.POST("/", SendMail)
	}

	return r
}
