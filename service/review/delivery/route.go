package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	review := r.Group("/api/reviews")
	{
		review.GET("/:review_id", GetReviewByID)
		review.GET("/", GetPaginatedReview)
		review.POST("/", CreateReview)
		review.PUT("/:review_id", UpdateReview)
		review.DELETE("/:review_id", DeleteReview)
	}

	return r
}
