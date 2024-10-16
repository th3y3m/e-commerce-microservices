package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetDiscountResponse struct {
	DiscountID    int64   `json:"discount_id"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	IsDeleted     bool    `json:"is_deleted"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
type ProductDiscount struct {
	ProductID  int64 `json:"product_id"`
	DiscountID int64 `json:"discount_id"`
}
type GetProductDiscountsRequest struct {
	ProductID  *int64 `json:"product_id"`
	DiscountID *int64 `json:"discount_id"`
}
type GetProductPriceAfterDiscount struct {
	ProductID int64 `json:"product_id"`
}
type GetProductRequest struct {
	ProductID int64 `json:"product_id"`
}
type GetProductsRequest struct {
	SellerID    *int64      `json:"seller_id"`
	ProductName string      `json:"product_name"`
	Description string      `json:"description"`
	MinPrice    *float64    `json:"min_price"`
	MaxPrice    *float64    `json:"max_price"`
	MinQuantity *int        `json:"min_quantity"`
	MaxQuantity *int        `json:"max_quantity"`
	CategoryID  *int64      `json:"category_id"`
	ImageURL    string      `json:"image_url"`
	FromDate    time.Time   `json:"from_date"`
	ToDate      time.Time   `json:"to_date"`
	IsDeleted   *bool       `json:"is_deleted"`
	Paging      util.Paging `json:"paging"`
}
type DeleteProductRequest struct {
	ProductID int64 `json:"product_id"`
}

type GetProductResponse struct {
	ProductID   int64   `json:"product_id"`
	SellerID    int64   `json:"seller_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  int64   `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	IsDeleted   bool    `json:"is_deleted"`
}

type GetProductListResponse struct {
	ProductID       int64   `json:"product_id"`
	SellerID        int64   `json:"seller_id"`
	ProductName     string  `json:"product_name"`
	Description     string  `json:"description"`
	OriginalPrice   float64 `json:"original_price"`
	DiscountedPrice float64 `json:"discounted_price"`
	Quantity        int     `json:"quantity"`
	CategoryID      int64   `json:"category_id"`
	ImageURL        string  `json:"image_url"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	IsDeleted       bool    `json:"is_deleted"`
}

type CreateProductRequest struct {
	SellerID    int64   `json:"seller_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  int64   `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}
type UpdateProductRequest struct {
	ProductID   int64   `json:"product_id"`
	SellerID    int64   `json:"seller_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  int64   `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}
