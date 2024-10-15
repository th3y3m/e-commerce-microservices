package model

import "time"

type GetDiscountRequest struct {
	DiscountID int64 `json:"discount_id"`
}

type DeleteDiscountRequest struct {
	DiscountID int64 `json:"discount_id"`
}

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

type CreateDiscountRequest struct {
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

type UpdateDiscountRequest struct {
	DiscountID    int64     `json:"discount_id"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsDeleted     bool      `json:"is_deleted"`
}
