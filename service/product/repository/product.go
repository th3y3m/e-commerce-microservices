package repository

import "time"

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
