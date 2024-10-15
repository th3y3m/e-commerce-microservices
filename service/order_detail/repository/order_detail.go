package repository

// OrderDetail represents a orderDetail in the system
type OrderDetail struct {
	OrderID   int64   `gorm:"primaryKey;column:order_id"`
	ProductID int64   `gorm:"primaryKey;column:product_id"`
	Quantity  int     `gorm:"column:quantity"`
	UnitPrice float64 `gorm:"column:unit_price"`
}
