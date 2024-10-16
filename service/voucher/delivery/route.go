package delivery

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	voucher := r.Group("/api/vouchers")
	{
		voucher.GET("/:voucher_id", GetVoucherByID)
		voucher.GET("", GetPaginatedVoucher)
		voucher.POST("", CreateVoucher)
		voucher.PUT("", UpdateVoucher)
		voucher.DELETE("", DeleteVoucher)
	}

	return r
}
