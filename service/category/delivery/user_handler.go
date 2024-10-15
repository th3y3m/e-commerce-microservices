package delivery

import (
	"th3y3m/e-commerce-microservices/service/category/dependency_injection"
	"th3y3m/e-commerce-microservices/service/category/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCategoryByID(c *gin.Context) {
	var req model.GetCategoryRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCategoryUsecaseProvider()

	category, err := module.GetCategory(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, category)
}

func GetAllCategorys(c *gin.Context) {
	module := dependency_injection.NewCategoryUsecaseProvider()

	categorys, err := module.GetAllCategorys(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, categorys)
}

func CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCategoryUsecaseProvider()

	category, err := module.CreateCategory(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, category)
}

func UpdateCategory(c *gin.Context) {
	var req model.UpdateCategoryRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCategoryUsecaseProvider()

	category, err := module.UpdateCategory(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, category)
}

func DeleteCategory(c *gin.Context) {
	var req model.DeleteCategoryRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewCategoryUsecaseProvider()

	err = module.DeleteCategory(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Category deleted successfully",
	})
}
