package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	user := r.Group("/api/users")
	{
		user.GET("/get-user", GetUser)
		user.GET("", GetPaginatedUser)
		user.POST("", CreateUser)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
	}

	return r
}
