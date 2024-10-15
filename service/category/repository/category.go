package repository

import "time"

// Category represents a category in the system
type Category struct {
	CategoryID   int64     `gorm:"primaryKey;column:category_id;autoIncrement"`
	CategoryName string    `gorm:"column:category_name"`
	IsDeleted    bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
