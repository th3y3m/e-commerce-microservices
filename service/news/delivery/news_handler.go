package delivery

import (
	"strconv"
	"th3y3m/e-commerce-microservices/service/news/dependency_injection"
	"th3y3m/e-commerce-microservices/service/news/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetNewsByID(c *gin.Context) {
	newID := c.Param("new_id")
	module := dependency_injection.NewNewsUsecaseProvider()

	id, err := strconv.ParseInt(newID, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	var req model.GetNewRequest
	req.NewsID = id

	new, err := module.GetNews(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, new)
}

func GetAllNews(c *gin.Context) {
	module := dependency_injection.NewNewsUsecaseProvider()

	news, err := module.GetAllNews(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, news)
}

func CreateNews(c *gin.Context) {
	var req model.CreateNewsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewNewsUsecaseProvider()

	new, err := module.CreateNews(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, new)
}

func UpdateNews(c *gin.Context) {
	var req model.UpdateNewsRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewNewsUsecaseProvider()

	new, err := module.UpdateNews(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, new)
}

func DeleteNews(c *gin.Context) {
	newID := c.Param("new_id")
	module := dependency_injection.NewNewsUsecaseProvider()

	id, err := strconv.ParseInt(newID, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	var req model.DeleteNewsRequest
	req.NewsID = id

	err = module.DeleteNews(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "New deleted successfully",
	})
}

func GetPaginatedNews(c *gin.Context) {
	module := dependency_injection.NewNewsUsecaseProvider()

	var req model.GetNewsRequest
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

	news, err := module.GetNewsList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, news)
}
