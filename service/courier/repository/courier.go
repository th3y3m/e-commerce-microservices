package repository

import "time"

// Courier represents a courier in the system
type Courier struct {
	CourierID   int64     `gorm:"primaryKey;column:courier_id;autoIncrement"`
	CourierName string    `gorm:"column:courier_name"`
	IsDeleted   bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt   time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
