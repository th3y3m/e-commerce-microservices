package delivery

import (
	"th3y3m/e-commerce-microservices/service/voucher/dependency_injection"
	"th3y3m/e-commerce-microservices/service/voucher/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetVoucherByID(c *gin.Context) {
	module := dependency_injection.NewVoucherUsecaseProvider()

	var req model.GetVoucherRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	voucher, err := module.GetVoucher(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, voucher)
}

func GetAllVouchers(c *gin.Context) {
	module := dependency_injection.NewVoucherUsecaseProvider()

	vouchers, err := module.GetAllVouchers(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, vouchers)
}

func CreateVoucher(c *gin.Context) {
	var req model.CreateVoucherRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewVoucherUsecaseProvider()

	voucher, err := module.CreateVoucher(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, voucher)
}

func UpdateVoucher(c *gin.Context) {
	var req model.UpdateVoucherRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	module := dependency_injection.NewVoucherUsecaseProvider()

	voucher, err := module.UpdateVoucher(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, voucher)
}

func DeleteVoucher(c *gin.Context) {
	module := dependency_injection.NewVoucherUsecaseProvider()

	var req model.DeleteVoucherRequest
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err = module.DeleteVoucher(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Voucher deleted successfully",
	})
}

func GetPaginatedVoucher(c *gin.Context) {
	module := dependency_injection.NewVoucherUsecaseProvider()

	var req model.GetVouchersRequest
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
	vouchers, err := module.GetVoucherList(c, &req)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, vouchers)
}
