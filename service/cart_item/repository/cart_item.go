package repository

// CartItem represents a cartItem in the system
type CartItem struct {
	CartID    int64 `gorm:"primaryKey;column:cart_id;autoIncrement"`
	ProductID int64 `gorm:"primaryKey;column:product_id"`
	Quantity  int   `gorm:"column:quantity"`
}
