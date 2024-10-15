package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	category := r.Group("/api/categories")
	{
		category.GET("/:category_id", GetCategoryByID)
		category.POST("/", CreateCategory)
		category.PUT("/", UpdateCategory)
		category.DELETE("/", DeleteCategory)
	}

	return r
}
