package repository

// ProductDiscount represents a productDiscount in the system
type ProductDiscount struct {
	ProductID  int64 `gorm:"primaryKey;column:product_id"`
	DiscountID int64 `gorm:"primaryKey;column:discount_id"`
}
