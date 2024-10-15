package model

type GetOrderDetailRequest struct {
	OrderID   int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
}

type GetOrderDetailsRequest struct {
	OrderID   *int64 `json:"order_id"`
	ProductID *int64 `json:"product_id"`
}

type DeleteOrderDetailRequest struct {
	OrderID   int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
}

type GetOrderDetailResponse struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type CreateOrderDetailRequest struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type UpdateOrderDetailRequest struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
