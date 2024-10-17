package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type Order struct {
	OrderID               int64     `json:"order_id"`
	CustomerID            int64     `json:"customer_id"`
	OrderDate             string    `json:"order_date"`
	TotalAmount           float64   `json:"total_amount"`
	OrderStatus           string    `json:"order_status"`
	ShippingAddress       string    `json:"shipping_address"`
	CourierID             int64     `json:"courier_id"`
	FreightPrice          float64   `json:"freight_price"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	ActualDeliveryDate    time.Time `json:"actual_delivery_date"`
	VoucherID             int64     `json:"voucher_id"`
	IsDeleted             bool      `json:"is_deleted"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
type GetVoucherRequest struct {
	VoucherID int64 `json:"voucher_id"`
}

type GetVouchersRequest struct {
	DiscountType string      `json:"discount_type"`
	IsAvailable  *bool       `json:"is_available"`
	IsDeleted    *bool       `json:"is_deleted"`
	Paging       util.Paging `json:"paging"`
}

type CheckVoucherUsageRequest struct {
	VoucherID int64 `json:"voucher_id"`
	Order     Order `json:"order"`
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
