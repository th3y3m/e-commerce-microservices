package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type PlaceOrderRequest struct {
	UserId        int64   `json:"user_id"`
	CartId        int64   `json:"cart_id"`
	CourierID     int64   `json:"courier_id"`
	VoucherID     int64   `json:"voucher_id"`
	PaymentMethod string  `json:"payment_method"`
	ShipAddress   string  `json:"ship_address"`
	Freight       float64 `json:"freight"`
}

type SendOrderDetailsRequest struct {
	Customer     User          `json:"user" binding:"required"`
	Order        Order         `json:"order" binding:"required"`
	OrderDetails []OrderDetail `json:"order_details" binding:"required"`
}
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
	Provider     string `json:"provider"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Token        string `json:"token"`
	TokenExpires string `json:"token_expires"`
	IsVerified   bool   `json:"is_verified"`
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
type CheckVoucherUsageRequest struct {
	VoucherID int64 `json:"voucher_id"`
	Order     Order `json:"order"`
}

type CheckVoucherUsageResponse struct {
	Valid bool `json:"valid"`
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
type User struct {
	UserID       int64     `json:"user_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Address      string    `json:"address"`
	Role         string    `json:"role"`
	ImageURL     string    `json:"image_url"`
	Provider     string    `json:"provider"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"token_expires"`
	IsVerified   bool      `json:"is_verified"`
	IsDeleted    bool      `json:"is_deleted"`
}

type OrderDetail struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
type Order struct {
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
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
