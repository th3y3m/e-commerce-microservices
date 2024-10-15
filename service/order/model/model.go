package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetOrderRequest struct {
	OrderID int64 `json:"order_id"`
}

type GetOrdersRequest struct {
	CustomerID            *int64      `json:"customer_id"`
	OrderDate             time.Time   `json:"order_date"`
	MinAmount             *float64    `json:"min_amount"`
	MaxAmount             *float64    `json:"max_amount"`
	OrderStatus           string      `json:"order_status"`
	ShippingAddress       string      `json:"shipping_address"`
	CourierID             *int64      `json:"courier_id"`
	FreightPrice          *float64    `json:"freight_price"`
	EstimatedDeliveryDate time.Time   `json:"estimated_delivery_date"`
	ActualDeliveryDate    time.Time   `json:"actual_delivery_date"`
	VoucherID             *int64      `json:"voucher_id"`
	FromDate              time.Time   `json:"from_date"`
	ToDate                time.Time   `json:"to_date"`
	IsDeleted             *bool       `json:"is_deleted"`
	Paging                util.Paging `json:"paging"`
}

type DeleteOrderRequest struct {
	OrderID int64 `json:"order_id"`
}

type GetOrderResponse struct {
	OrderID               int64   `json:"order_id"`
	CustomerID            int64   `json:"customer_id"`
	OrderDate             string  `json:"order_date"`
	TotalAmount           float64 `json:"total_amount"`
	OrderStatus           string  `json:"order_status"`
	ShippingAddress       string  `json:"shipping_address"`
	CourierID             int64   `json:"courier_id"`
	FreightPrice          float64 `json:"freight_price"`
	EstimatedDeliveryDate string  `json:"estimated_delivery_date"`
	ActualDeliveryDate    string  `json:"actual_delivery_date"`
	VoucherID             int64   `json:"voucher_id"`
	IsDeleted             bool    `json:"is_deleted"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
}

type CreateOrderRequest struct {
	CustomerID            int64     `json:"customer_id"`
	TotalAmount           float64   `json:"total_amount"`
	OrderStatus           string    `json:"order_status"`
	ShippingAddress       string    `json:"shipping_address"`
	CourierID             int64     `json:"courier_id"`
	FreightPrice          float64   `json:"freight_price"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	ActualDeliveryDate    time.Time `json:"actual_delivery_date"`
	VoucherID             int64     `json:"voucher_id"`
}

type UpdateOrderRequest struct {
	OrderID               int64     `json:"order_id"`
	CustomerID            int64     `json:"customer_id"`
	OrderDate             time.Time `json:"order_date"`
	TotalAmount           float64   `json:"total_amount"`
	OrderStatus           string    `json:"order_status"`
	ShippingAddress       string    `json:"shipping_address"`
	CourierID             int64     `json:"courier_id"`
	FreightPrice          float64   `json:"freight_price"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	ActualDeliveryDate    time.Time `json:"actual_delivery_date"`
	VoucherID             int64     `json:"voucher_id"`
	IsDeleted             bool      `json:"is_deleted"`
}
