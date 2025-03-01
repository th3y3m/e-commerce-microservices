package model

import "time"

type UpdateUserRequest struct {
	UserID       int64     `json:"user_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Address      string    `json:"address"`
	Role         string    `json:"role"`
	ImageURL     string    `json:"image_url"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"token_expires"`
	IsVerified   bool      `json:"is_verified"`
	IsDeleted    bool      `json:"is_deleted"`
}
type CreateUserRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	ImageURL   string `json:"image_url"`
	Provider   string `json:"provider"`
	IsVerified bool   `json:"is_verified"`
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
	Provider     string    `gorm:"column:provider"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	Token        string    `gorm:"column:token"`
	TokenExpires time.Time `gorm:"column:token_expires"`
	IsVerified   bool      `gorm:"column:is_verified;default:false"`
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
