package delivery

import (
	"th3y3m/e-commerce-microservices/service/review/dependency_injection"
	"th3y3m/e-commerce-microservices/service/review/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetReviewByID(c *gin.Context) {
	module := dependency_injection.NewReviewUsecaseProvider()

	var req model.GetReviewRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	review, err := module.GetReview(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, review)
}

func GetAllReviews(c *gin.Context) {
	module := dependency_injection.NewReviewUsecaseProvider()

	reviews, err := module.GetAllReviews(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, reviews)
}

func CreateReview(c *gin.Context) {
	var req model.CreateReviewRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewReviewUsecaseProvider()

	review, err := module.CreateReview(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, review)
}

func UpdateReview(c *gin.Context) {
	var req model.UpdateReviewRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewReviewUsecaseProvider()

	review, err := module.UpdateReview(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, review)
}

func DeleteReview(c *gin.Context) {
	module := dependency_injection.NewReviewUsecaseProvider()

	var req model.DeleteReviewRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteReview(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Review deleted successfully",
	})
}

func GetPaginatedReview(c *gin.Context) {
	module := dependency_injection.NewReviewUsecaseProvider()

	var req model.GetReviewsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
	if req.Paging.PageIndex == 0 {
		req.Paging.PageIndex = 1
	}
	if req.Paging.PageSize == 0 {
		req.Paging.PageSize = 10
	}
	reviews, err := module.GetReviewList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, reviews)
}
