package repository

import "time"

// New represents a new in the system
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
