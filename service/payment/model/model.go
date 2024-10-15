package model

import (
	"th3y3m/e-commerce-microservices/pkg/util"
	"time"
)

type GetPaymentRequest struct {
	PaymentID int64 `json:"payment_id"`
}

type GetPaymentsRequest struct {
	OrderID       *int64      `json:"order_id"`
	MinAmount     *float64    `json:"min_amount"`
	MaxAmount     *float64    `json:"max_amount"`
	FromDate      time.Time   `json:"from_date"`
	ToDate        time.Time   `json:"to_date"`
	PaymentMethod string      `json:"payment_method"`
	PaymentStatus string      `json:"payment_status"`
	Paging        util.Paging `json:"paging"`
}

type DeletePaymentRequest struct {
	PaymentID int64 `json:"payment_id"`
}

type GetPaymentResponse struct {
	PaymentID        int64   `json:"payment_id"`
	OrderID          int64   `json:"order_id"`
	PaymentAmount    float64 `json:"payment_amount"`
	PaymentDate      string  `json:"payment_date"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentStatus    string  `json:"payment_status"`
	PaymentSignature string  `json:"payment_signature"`
}

type CreatePaymentRequest struct {
	OrderID          int64   `json:"order_id"`
	PaymentAmount    float64 `json:"payment_amount"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentStatus    string  `json:"payment_status"`
	PaymentSignature string  `json:"payment_signature"`
}

type UpdatePaymentRequest struct {
	PaymentID        int64   `json:"payment_id"`
	PaymentAmount    float64 `json:"payment_amount"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentStatus    string  `json:"payment_status"`
	PaymentSignature string  `json:"payment_signature"`
}
