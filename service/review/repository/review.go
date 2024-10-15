package repository

import "time"

// Review represents a review in the system
type Review struct {
	ReviewID  int64     `gorm:"primaryKey;column:review_id;autoIncrement"`
	ProductID int64     `gorm:"column:product_id"`
	UserID    int64     `gorm:"column:user_id"`
	Rating    int       `gorm:"column:rating"`
	Comment   string    `gorm:"column:comment"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
