package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	user := r.Group("/api/users")
	{
		user.GET("/:user_id", GetUserByID)
		user.GET("/", GetPaginatedUser)
		user.POST("/", CreateUser)
		user.PUT("/:user_id", UpdateUser)
		user.DELETE("/:user_id", DeleteUser)
	}

	return r
}
