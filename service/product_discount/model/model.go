package model

type GetProductDiscountRequest struct {
	ProductID  int64 `json:"product_id"`
	DiscountID int64 `json:"discount_id"`
}

type GetProductDiscountsRequest struct {
	ProductID  *int64 `json:"product_id"`
	DiscountID *int64 `json:"discount_id"`
}

type DeleteProductDiscountRequest struct {
	ProductID  int64 `json:"product_id"`
	DiscountID int64 `json:"discount_id"`
}

type GetProductDiscountResponse struct {
	ProductID  int64 `json:"product_id"`
	DiscountID int64 `json:"discount_id"`
}

type CreateProductDiscountRequest struct {
	ProductID  int64 `json:"product_id"`
	DiscountID int64 `json:"discount_id"`
}

type UpdateProductDiscountRequest struct {
	ProductID  int64 `json:"product_id"`
	DiscountID int64 `json:"discount_id"`
}
