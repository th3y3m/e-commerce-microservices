package model

import "time"

type CreatePaymentRequest struct {
	OrderID          int64   `json:"order_id"`
	PaymentAmount    float64 `json:"payment_amount"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentStatus    string  `json:"payment_status"`
	PaymentSignature string  `json:"payment_signature"`
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
type Payment struct {
	PaymentID        int64     `gorm:"primaryKey;column:payment_id;autoIncrement"`
	OrderID          int64     `gorm:"column:order_id"`
	PaymentAmount    float64   `gorm:"column:payment_amount"`
	PaymentDate      time.Time `gorm:"autoCreateTime;column:payment_date"`
	PaymentMethod    string    `gorm:"column:payment_method"`
	PaymentStatus    string    `gorm:"column:payment_status"`
	PaymentSignature string    `gorm:"column:payment_signature"`
}
type PaymentResponse struct {
	IsSuccessful bool   `json:"is_successful"`
	RedirectUrl  string `json:"redirect_url"`
}

type User struct {
	UserID       int64     `gorm:"primaryKey;autoIncrement;column:user_id"`
	Email        string    `gorm:"unique;not null;column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	FullName     string    `gorm:"column:full_name"`
	PhoneNumber  string    `gorm:"column:phone_number"`
	Address      string    `gorm:"column:address"`
	Role         string    `gorm:"column:role"`
	ImageURL     string    `gorm:"column:image_url"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	Token        string    `gorm:"column:token"`
	TokenExpires time.Time `gorm:"column:token_expires"`
	IsVerified   bool      `gorm:"column:is_verified"`
	IsDeleted    bool      `gorm:"column:is_deleted;default:false"`
}

type OrderDetail struct {
	OrderID   int64   `gorm:"primaryKey;column:order_id"`
	ProductID int64   `gorm:"primaryKey;column:product_id"`
	Quantity  int     `gorm:"column:quantity"`
	UnitPrice float64 `gorm:"column:unit_price"`
}
type Order struct {
	OrderID               int64     `gorm:"primaryKey;column:order_id;autoIncrement"`
	CustomerID            int64     `gorm:"column:customer_id"`
	OrderDate             time.Time `gorm:"autoCreateTime;column:order_date"`
	TotalAmount           float64   `gorm:"column:total_amount"`
	OrderStatus           string    `gorm:"column:order_status"`
	ShippingAddress       string    `gorm:"column:shipping_address"`
	CourierID             int64     `gorm:"column:courier_id"`
	FreightPrice          float64   `gorm:"column:freight_price"`
	EstimatedDeliveryDate time.Time `gorm:"column:estimated_delivery_date"`
	ActualDeliveryDate    time.Time `gorm:"column:actual_delivery_date"`
	VoucherID             int64     `gorm:"column:voucher_id"`
	IsDeleted             bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt             time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt             time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}

type Product struct {
	ProductID   int64     `gorm:"primaryKey;column:product_id;autoIncrement"`
	SellerID    int64     `gorm:"column:seller_id"`
	ProductName string    `gorm:"column:product_name"`
	Description string    `gorm:"column:description"`
	Price       float64   `gorm:"column:price"`
	Quantity    int       `gorm:"column:quantity"`
	CategoryID  int64     `gorm:"column:category_id"`
	ImageURL    string    `gorm:"column:image_url"`
	CreatedAt   time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	IsDeleted   bool      `gorm:"column:is_deleted;default:false"`
}
