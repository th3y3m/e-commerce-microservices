package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetUserRequest struct {
	UserID *int64 `json:"user_id"`
	Email  string `json:"email"`
}

type GetUsersRequest struct {
	Email       string      `json:"email"`
	FullName    string      `json:"full_name"`
	PhoneNumber string      `json:"phone_number"`
	Address     string      `json:"address"`
	Role        string      `json:"role"`
	FromDate    time.Time   `json:"from_date"`
	ToDate      time.Time   `json:"to_date"`
	IsDeleted   *bool       `json:"is_deleted"`
	Paging      util.Paging `json:"paging"`
}

type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserResponse struct {
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
	Role         string `json:"role"`
	ImageURL     string `json:"image_url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Token        string `json:"token"`
	TokenExpires string `json:"token_expires"`
	IsVerified   bool   `json:"is_verified"`
	IsDeleted    bool   `json:"is_deleted"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	UserID       int64     `json:"user_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Address      string    `json:"address"`
	Role         string    `json:"role"`
	ImageURL     string    `json:"image_url"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"token_expires"`
	IsVerified   bool      `json:"is_verified"`
	IsDeleted    bool      `json:"is_deleted"`
}
