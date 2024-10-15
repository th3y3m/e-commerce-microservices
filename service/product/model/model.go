package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

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
