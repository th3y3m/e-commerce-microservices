package model

import "time"

type User struct {
	UserID       int64     `json:"user_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Address      string    `json:"address"`
	Role         string    `json:"role"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"token_expires"`
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

type Product struct {
	ProductID   int64     `json:"product_id"`
	SellerID    int64     `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CategoryID  int64     `json:"category_id"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDeleted   bool      `json:"is_deleted"`
}

type SendOrderDetailsRequest struct {
	Customer     User          `json:"customer" binding:"required"`
	Order        Order         `json:"order" binding:"required"`
	OrderDetails []OrderDetail `json:"order_details" binding:"required"`
}
