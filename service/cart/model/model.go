package model

import (
	"time"
)

type Item struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
type DeleteCartItemRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
}
type UpdateOrCreateRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
type GetCartItemResponse struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
type GetCartItemsRequest struct {
	CartID    *int64 `json:"cart_id"`
	ProductID *int64 `json:"product_id"`
}
type CartItem struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type SeachQueryRequest struct {
	UserID    int64 `json:"user_id"`
	IsDeleted bool  `json:"is_deleted"`
}

type GetCartRequest struct {
	CartID int64 `json:"cart_id"`
}

type GetCartsRequest struct {
	CartID    int64     `json:"cart_id"`
	UserID    int64     `json:"user_id"`
	IsDeleted bool      `json:"is_deleted"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
}

type DeleteCartRequest struct {
	CartID int64 `json:"cart_id"`
}

type GetCartResponse struct {
	CartID    int64  `json:"cart_id"`
	UserID    int64  `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateCartRequest struct {
	CartID int64 `json:"cart_id"`
	UserID int64 `json:"user_id"`
}

type UpdateCartRequest struct {
	CartID    int64 `json:"cart_id"`
	IsDeleted bool  `json:"is_deleted"`
}
