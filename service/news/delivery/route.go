package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	new := r.Group("/api/news")
	{
		new.GET("/:new_id", GetNewsByID)
		new.GET("", GetPaginatedNews)
		new.POST("", CreateNews)
		new.PUT("/:new_id", UpdateNews)
		new.DELETE("/:new_id", DeleteNews)
	}

	return r
}
