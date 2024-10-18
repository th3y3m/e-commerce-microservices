package delivery

import (
	"strconv"
	"th3y3m/e-commerce-microservices/service/product/dependency_injection"
	"th3y3m/e-commerce-microservices/service/product/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetProductByID(c *gin.Context) {
	module := dependency_injection.NewProductUsecaseProvider()

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	req := model.GetProductRequest{
		ProductID: productID,
	}

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
	module := dependency_injection.NewProductUsecaseProvider()

	var req model.DeleteProductRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
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
	module := dependency_injection.NewProductUsecaseProvider()

	var req model.GetProductsRequest
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

func GetProductPriceAfterDiscount(c *gin.Context) {
	module := dependency_injection.NewProductUsecaseProvider()

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	req := model.GetProductPriceAfterDiscount{
		ProductID: productID,
	}

	price, err := module.GetProductPriceAfterDiscount(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, price)
}
