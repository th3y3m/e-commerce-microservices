package repository

import "time"

// Discount represents a discount in the system
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
