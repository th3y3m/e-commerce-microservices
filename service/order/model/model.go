package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetOrderDetailsRequest struct {
	OrderID   *int64 `json:"order_id"`
	ProductID *int64 `json:"product_id"`
}
type GetUserResponse struct {
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
	Role         string `json:"role"`
	ImageURL     string `json:"image_url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Token        string `json:"token"`
	TokenExpires string `json:"token_expires"`
	IsDeleted    bool   `json:"is_deleted"`
}
type GetUserRequest struct {
	UserID *int64 `json:"user_id"`
	Email  string `json:"email"`
}
type UpdateProductRequest struct {
	ProductID   int64   `json:"product_id"`
	SellerID    int64   `json:"seller_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  int64   `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}
type CreatePaymentRequest struct {
	OrderID          int64   `json:"order_id"`
	PaymentAmount    float64 `json:"payment_amount"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentStatus    string  `json:"payment_status"`
	PaymentSignature string  `json:"payment_signature"`
}
type GetOrderDetailResponse struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
type CreateOrderDetailRequest struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type GetProductResponse struct {
	ProductID   int64   `json:"product_id"`
	SellerID    int64   `json:"seller_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  int64   `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	IsDeleted   bool    `json:"is_deleted"`
}
type GetProductRequest struct {
	ProductID int64 `json:"product_id"`
}
type GetCartItemsRequest struct {
	CartID    *int64 `json:"cart_id"`
	ProductID *int64 `json:"product_id"`
}

type GetCartItemResponse struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

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
