package model

import "th3y3m/e-commerce-microservices/pkg/util"

type GetReviewRequest struct {
	ReviewID int64 `json:"review_id"`
}

type GetReviewsRequest struct {
	ProductID *int64      `json:"product_id"`
	UserID    *int64      `json:"user_id"`
	MinRating *int        `json:"min_rating"`
	MaxRating *int        `json:"max_rating"`
	Comment   string      `json:"comment"`
	IsDeleted *bool       `json:"is_deleted"`
	Paging    util.Paging `json:"paging"`
}

type DeleteReviewRequest struct {
	ReviewID int64 `json:"review_id"`
}

type GetReviewResponse struct {
	ReviewID  int64  `json:"review_id"`
	ProductID int64  `json:"product_id"`
	UserID    int64  `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	IsDeleted bool   `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateReviewRequest struct {
	ProductID int64  `json:"product_id"`
	UserID    int64  `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

type UpdateReviewRequest struct {
	ReviewID  int64  `json:"review_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	IsDeleted bool   `json:"is_deleted"`
}
