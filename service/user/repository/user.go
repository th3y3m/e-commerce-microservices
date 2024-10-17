package repository

import "time"

// User represents a user in the system
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
