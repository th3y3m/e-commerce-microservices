package util

import "math"

// PaginatedList represents a paginated list of items
type PaginatedList[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
	PageIndex  int `json:"page_index"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

func (p *PaginatedList[any]) GetTotalPages() {
	p.TotalPages = int(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
}

type Paging struct {
	PageIndex     int    `json:"page_index"`
	PageSize      int    `json:"page_size"`
	Sort          string `json:"sort"`
	SortDirection string `json:"sort_direction"`
}
