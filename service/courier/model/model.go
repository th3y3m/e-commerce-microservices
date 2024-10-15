package model

type GetCourierRequest struct {
	CourierID int64 `json:"courier_id"`
}

type DeleteCourierRequest struct {
	CourierID int64 `json:"courier_id"`
}

type GetCourierResponse struct {
	CourierID   int64  `json:"courier_id"`
	CourierName string `json:"courier_name"`
	IsDeleted   bool   `json:"is_deleted"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateCourierRequest struct {
	CourierName string `json:"courier_name"`
}

type UpdateCourierRequest struct {
	CourierID   int64  `json:"courier_id"`
	CourierName string `json:"courier_name"`
	IsDeleted   bool   `json:"is_deleted"`
}
