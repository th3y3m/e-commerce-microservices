package repository

import "time"

// Cart represents a cart in the system
type Cart struct {
	CartID    int64     `gorm:"primaryKey;column:cart_id;autoIncrement"`
	UserID    int64     `gorm:"column:user_id"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
