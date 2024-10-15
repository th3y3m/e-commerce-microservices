package repository

import "time"

// Payment represents a payment in the system
type Payment struct {
	PaymentID        int64     `gorm:"primaryKey;column:payment_id;autoIncrement"`
	OrderID          int64     `gorm:"column:order_id"`
	PaymentAmount    float64   `gorm:"column:payment_amount"`
	PaymentDate      time.Time `gorm:"autoCreateTime;column:payment_date"`
	PaymentMethod    string    `gorm:"column:payment_method"`
	PaymentStatus    string    `gorm:"column:payment_status"`
	PaymentSignature string    `gorm:"column:payment_signature"`
}
