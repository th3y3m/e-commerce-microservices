package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Voucher struct {
	VoucherID          int64     `gorm:"primaryKey;column:voucher_id;autoIncrement"`
	VoucherCode        string    `gorm:"unique;not null;column:voucher_code"`
	DiscountType       string    `gorm:"column:discount_type"`
	DiscountValue      float64   `gorm:"column:discount_value"`
	MinimumOrderAmount float64   `gorm:"column:minimum_order_amount"`
	MaxDiscountAmount  float64   `gorm:"column:max_discount_amount"`
	StartDate          time.Time `gorm:"column:start_date"`
	EndDate            time.Time `gorm:"column:end_date"`
	UsageLimit         int       `gorm:"column:usage_limit"`
	UsageCount         int       `gorm:"column:usage_count"`
	IsDeleted          bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt          time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt          time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
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

type Review struct {
	ReviewID  int64     `gorm:"primaryKey;column:review_id;autoIncrement"`
	ProductID int64     `gorm:"column:product_id"`
	UserID    int64     `gorm:"column:user_id"`
	Rating    int       `gorm:"column:rating"`
	Comment   string    `gorm:"column:comment"`
	CreatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
}
type ProductDiscount struct {
	ProductID  int64 `gorm:"primaryKey;column:product_id"`
	DiscountID int64 `gorm:"primaryKey;column:discount_id"`
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

type FreightRate struct {
	FreightRateID int64     `gorm:"primaryKey;column:rate_id;autoIncrement"`
	CourierID     int64     `gorm:"column:courier_id"`
	DistanceMinKM float64   `gorm:"column:distance_min_km"`
	DistanceMaxKM float64   `gorm:"column:distance_max_km"`
	CostPerKM     float64   `gorm:"column:cost_per_km"`
	IsDeleted     bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
type Discount struct {
	DiscountID    int64     `gorm:"primaryKey;autoIncrement;column:discount_id"`
	DiscountType  string    `gorm:"column:discount_type"`
	DiscountValue float64   `gorm:"column:discount_value"`
	StartDate     time.Time `gorm:"column:start_date"`
	EndDate       time.Time `gorm:"column:end_date"`
	IsDeleted     bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
type Courier struct {
	CourierID   int64     `gorm:"primaryKey;column:courier_id;autoIncrement"`
	CourierName string    `gorm:"column:courier_name"`
	IsDeleted   bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt   time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
type Category struct {
	CategoryID   int64     `gorm:"primaryKey;column:category_id;autoIncrement"`
	CategoryName string    `gorm:"column:category_name"`
	IsDeleted    bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
type Cart struct {
	CartID    int64     `gorm:"primaryKey;column:cart_id;autoIncrement"`
	UserID    int64     `gorm:"column:user_id"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
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

type CartItem struct {
	CartID    int64 `gorm:"primaryKey;column:cart_id"`
	ProductID int64 `gorm:"primaryKey;column:product_id"`
	Quantity  int   `gorm:"column:quantity"`
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
	IsVerified   bool      `gorm:"column:is_verified;default:false"`
	IsDeleted    bool      `gorm:"column:is_deleted;default:false"`
}
type News struct {
	NewsID        int64     `gorm:"primaryKey;autoIncrement;column:news_id"`
	Title         string    `gorm:"column:title"`
	Content       string    `gorm:"column:content"`
	PublishedDate time.Time `gorm:"autoCreateTime;column:published_date"`
	AuthorID      int64     `gorm:"column:author_id"`
	ImageURL      string    `gorm:"column:image_url"`
	Category      string    `gorm:"column:category"`
	IsDeleted     bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}

func main() {

	dsn := "host=localhost user=postgres password=12345 dbname=ECommerceDb port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate the schemas
	db.AutoMigrate(&Product{}, &User{}, &Cart{}, &CartItem{}, &Category{}, &Courier{}, &Discount{}, &FreightRate{}, &Order{}, &OrderDetail{}, &ProductDiscount{}, &Review{}, &Payment{}, &Voucher{}, &News{})

	log.Println("Database migration completed successfully!")

}
