package model

type GetFreightRateRequest struct {
	FreightRateID int64 `json:"freightRate_id"`
}

type DeleteFreightRateRequest struct {
	FreightRateID int64 `json:"freightRate_id"`
}

type GetFreightRateResponse struct {
	FreightRateID int64   `gorm:"primaryKey;column:rate_id;autoIncrement"`
	CourierID     int64   `gorm:"column:courier_id"`
	DistanceMinKM float64 `gorm:"column:distance_min_km"`
	DistanceMaxKM float64 `gorm:"column:distance_max_km"`
	CostPerKM     float64 `gorm:"column:cost_per_km"`
	IsDeleted     bool    `gorm:"column:is_deleted;default:false"`
	CreatedAt     string  `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt     string  `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}

type CreateFreightRateRequest struct {
	CourierID     int64   `json:"courier_id"`
	DistanceMinKM float64 `json:"distance_min_km"`
	DistanceMaxKM float64 `json:"distance_max_km"`
	CostPerKM     float64 `json:"cost_per_km"`
}

type UpdateFreightRateRequest struct {
	FreightRateID int64   `json:"freightRate_id"`
	CourierID     int64   `json:"courier_id"`
	DistanceMinKM float64 `json:"distance_min_km"`
	DistanceMaxKM float64 `json:"distance_max_km"`
	CostPerKM     float64 `json:"cost_per_km"`
	IsDeleted     bool    `json:"is_deleted"`
}
