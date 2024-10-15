package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetNewRequest struct {
	NewsID int64 `json:"new_id"`
}

type GetNewsRequest struct {
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	FromDate  time.Time   `json:"from_date"`
	ToDate    time.Time   `json:"to_date"`
	AuthorID  *int64      `json:"author_id"`
	Category  string      `json:"category"`
	IsDeleted *bool       `json:"is_deleted"`
	Paging    util.Paging `json:"paging"`
}

type DeleteNewsRequest struct {
	NewsID int64 `json:"new_id"`
}

type GetNewsResponse struct {
	NewsID        int64  `json:"news_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	PublishedDate string `json:"published_date"`
	AuthorID      int64  `json:"author_id"`
	ImageURL      string `json:"image_url"`
	Category      string `json:"category"`
	IsDeleted     bool   `json:"is_deleted"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CreateNewsRequest struct {
	Title     string
	Content   string
	AuthorID  int64
	ImageURL  string
	Category  string
	IsDeleted bool
}

type UpdateNewsRequest struct {
	NewsID    int64
	Title     string
	Content   string
	AuthorID  int64
	ImageURL  string
	Category  string
	IsDeleted bool
}
