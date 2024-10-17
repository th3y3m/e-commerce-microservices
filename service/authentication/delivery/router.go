package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	authen := r.Group("/api/authen")
	{
		authen.POST("/login", Login)
		authen.POST("/register", Register)
		authen.GET("/verify-email", VerifyUserEmail)
	}

	return r
}
