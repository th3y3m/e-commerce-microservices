package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetVoucherRequest struct {
	VoucherID int64 `json:"voucher_id"`
}

type GetVouchersRequest struct {
	DiscountType string      `json:"discount_type"`
	IsAvailable  *bool       `json:"is_available"`
	IsDeleted    *bool       `json:"is_deleted"`
	Paging       util.Paging `json:"paging"`
}

type DeleteVoucherRequest struct {
	VoucherID int64 `json:"voucher_id"`
}

type GetVoucherResponse struct {
	VoucherID          int64   `json:"voucher_id"`
	VoucherCode        string  `json:"voucher_code"`
	DiscountType       string  `json:"discount_type"`
	DiscountValue      float64 `json:"discount_value"`
	MinimumOrderAmount float64 `json:"minimum_order_amount"`
	MaxDiscountAmount  float64 `json:"max_discount_amount"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	UsageLimit         int     `json:"usage_limit"`
	UsageCount         int     `json:"usage_count"`
	IsDeleted          bool    `json:"is_deleted"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type CreateVoucherRequest struct {
	VoucherCode        string    `json:"voucher_code"`
	DiscountType       string    `json:"discount_type"`
	DiscountValue      float64   `json:"discount_value"`
	MinimumOrderAmount float64   `json:"minimum_order_amount"`
	MaxDiscountAmount  float64   `json:"max_discount_amount"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	UsageLimit         int       `json:"usage_limit"`
	UsageCount         int       `json:"usage_count"`
}

type UpdateVoucherRequest struct {
	VoucherID          int64     `json:"voucher_id"`
	VoucherCode        string    `json:"voucher_code"`
	DiscountType       string    `json:"discount_type"`
	DiscountValue      float64   `json:"discount_value"`
	MinimumOrderAmount float64   `json:"minimum_order_amount"`
	MaxDiscountAmount  float64   `json:"max_discount_amount"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	UsageLimit         int       `json:"usage_limit"`
	UsageCount         int       `json:"usage_count"`
	IsDeleted          bool      `json:"is_deleted"`
}
