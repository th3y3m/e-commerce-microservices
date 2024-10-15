package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Base URL for the product service
const productServiceBaseURL = "http://localhost:8081/api/products"

// GetProductByID proxies the request to the product service
func GetProductByID(c *gin.Context) {
	productID := c.Param("product_id")
	resp, err := http.Get(productServiceBaseURL + "/" + productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// GetPaginatedProducts proxies the request to the product service
func GetPaginatedProducts(c *gin.Context) {
	page := c.DefaultQuery("page_index", "1")
	size := c.DefaultQuery("page_size", "10")
	resp, err := http.Get(productServiceBaseURL + "?page_index=" + page + "&page_size=" + size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// CreateProduct proxies the POST request to the product service
func CreateProduct(c *gin.Context) {
	resp, err := http.Post(productServiceBaseURL, "application/json", c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// UpdateProduct proxies the PUT request to the product service
func UpdateProduct(c *gin.Context) {
	productID := c.Param("product_id")
	req, err := http.NewRequest(http.MethodPut, productServiceBaseURL+"/"+productID, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// DeleteProduct proxies the DELETE request to the product service
func DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")
	req, err := http.NewRequest(http.MethodDelete, productServiceBaseURL+"/"+productID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reach product service"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
