package delivery

import (
	"strconv"
	"th3y3m/e-commerce-microservices/service/product/dependency_injection"
	"th3y3m/e-commerce-microservices/service/product/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetProductByID(c *gin.Context) {
	productID := c.Param("product_id")
	module := dependency_injection.NewProductUsecaseProvider()

	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	var req model.GetProductRequest
	req.ProductID = id

	product, err := module.GetProduct(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, product)
}

func GetAllProducts(c *gin.Context) {
	module := dependency_injection.NewProductUsecaseProvider()

	products, err := module.GetAllProducts(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, products)
}

func CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewProductUsecaseProvider()

	product, err := module.CreateProduct(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, product)
}

func UpdateProduct(c *gin.Context) {
	var req model.UpdateProductRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewProductUsecaseProvider()

	product, err := module.UpdateProduct(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, product)
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")
	module := dependency_injection.NewProductUsecaseProvider()

	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	var req model.DeleteProductRequest
	req.ProductID = id

	err = module.DeleteProduct(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product deleted successfully",
	})
}

func GetPaginatedProduct(c *gin.Context) {
	page := c.DefaultQuery("page_index", "1")
	limit := c.DefaultQuery("page_size", "10")
	module := dependency_injection.NewProductUsecaseProvider()

	p, err := strconv.Atoi(page)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	var req model.GetProductsRequest
	req.Paging.PageIndex = p
	req.Paging.PageSize = l

	products, err := module.GetProductList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, products)
}
