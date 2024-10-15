package repository

import "time"

// FreightRate represents a freightRate in the system
type FreightRate struct {
	FreightRateID int64     `gorm:"primaryKey;column:rate_id;autoIncrement"`
	CourierID     int64     `gorm:"column:courier_id"`
	DistanceMinKM float64   `gorm:"column:distance_min_km"`
	DistanceMaxKM float64   `gorm:"column:distance_max_km"`
	CostPerKM     float64   `gorm:"column:cost_per_km"`
	IsDeleted     bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp without time zone;column:created_at;default:current_timestamp"`
}
