package model

type GetCartItemRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
}

type GetCartItemsRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
}

type DeleteCartItemRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
}

type GetCartItemResponse struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CreateCartItemRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type UpdateCartItemRequest struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
