package repository

import "time"

// Order represents a order in the system
type Order struct {
	OrderID               int64     `gorm:"primaryKey;column:order_id;autoIncrement"`
	CustomerID            int64     `gorm:"column:customer_id"`
	OrderDate             time.Time `gorm:"autoCreateTime;column:order_date"`
	TotalAmount           float64   `gorm:"column:total_amount"`
	OrderStatus           string    `gorm:"column:order_status"`
	ShippingAddress       string    `gorm:"column:shipping_address"`
	CourierID             int64     `gorm:"column:courier_id"`
	FreightPrice          float64   `gorm:"column:freight_price"`
	EstimatedDeliveryDate time.Time `gorm:"column:estimated_delivery_date"`
	ActualDeliveryDate    time.Time `gorm:"column:actual_delivery_date"`
	VoucherID             int64     `gorm:"column:voucher_id"`
	IsDeleted             bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt             time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt             time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
