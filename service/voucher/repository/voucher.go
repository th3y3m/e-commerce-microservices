package repository

import "time"

// Voucher represents a voucher in the system
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
