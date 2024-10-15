package model

type GetCategoryRequest struct {
	CategoryID int64 `json:"category_id"`
}

type DeleteCategoryRequest struct {
	CategoryID int64 `json:"category_id"`
}

type GetCategoryResponse struct {
	CategoryID   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	IsDeleted    bool   `json:"is_deleted"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateCategoryRequest struct {
	CategoryName string `json:"category_name"`
}

type UpdateCategoryRequest struct {
	CategoryID   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	IsDeleted    bool   `json:"is_deleted"`
}
