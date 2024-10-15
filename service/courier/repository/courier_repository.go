package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type courierRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type ICourierRepository interface {
	Get(ctx context.Context, courierID int64) (*Courier, error)
	GetAll(ctx context.Context) ([]*Courier, error)
	Create(ctx context.Context, courier *Courier) (*Courier, error)
	Update(ctx context.Context, courier *Courier) (*Courier, error)
	Delete(ctx context.Context, courierID int64) error
}

func NewCourierRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) ICourierRepository {
	return &courierRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *courierRepository) Get(ctx context.Context, courierID int64) (*Courier, error) {
	pr.log.Infof("Fetching courier with ID: %d", courierID)
	cacheKey := fmt.Sprintf("courier:%d", courierID)

	// Try to get the courier from Redis cache
	cachedCourier, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var courier Courier
		if err := json.Unmarshal([]byte(cachedCourier), &courier); err == nil {
			pr.log.Infof("Courier found in cache: %+v", courier)
			return &courier, nil
		}
	}

	// If not found in cache, get from database
	var courier Courier
	result := pr.db.WithContext(ctx).First(&courier, courierID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching courier: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	courierJSON, _ := json.Marshal(courier)
	pr.redis.Set(ctx, cacheKey, courierJSON, 0)
	pr.log.Infof("Courier saved to cache: %d", courierID)

	return &courier, nil
}

func (pr *courierRepository) GetAll(ctx context.Context) ([]*Courier, error) {
	pr.log.Info("Fetching all couriers")
	cacheKey := "all_couriers"

	// Try to get the couriers from Redis cache
	cachedCouriers, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var couriers []*Courier
		if err := json.Unmarshal([]byte(cachedCouriers), &couriers); err == nil {
			pr.log.Infof("Couriers found in cache: %+v", couriers)
			return couriers, nil
		}
	}

	// If not found in cache, get from database
	var couriers []*Courier
	result := pr.db.WithContext(ctx).Find(&couriers)
	if result.Error != nil {
		pr.log.Errorf("Error fetching couriers: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	couriersJSON, _ := json.Marshal(couriers)
	pr.redis.Set(ctx, cacheKey, couriersJSON, 0)
	pr.log.Info("Couriers saved to cache")

	return couriers, nil
}

func (pr *courierRepository) Create(ctx context.Context, courier *Courier) (*Courier, error) {
	pr.log.Infof("Creating courier: %+v", courier)
	if err := pr.db.WithContext(ctx).Create(courier).Error; err != nil {
		pr.log.Errorf("Error creating courier: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("courier:%d", courier.CourierID)

	courierJSON, _ := json.Marshal(courier)
	pr.redis.Set(ctx, cacheKey, courierJSON, 0)
	pr.log.Infof("Courier saved to cache: %d", courier.CourierID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_couriers")
	pr.log.Info("Invalidated cache for all couriers")

	// Return the newly created courier (with any updated fields)
	return courier, nil
}

func (pr *courierRepository) Update(ctx context.Context, courier *Courier) (*Courier, error) {
	pr.log.Infof("Updating courier: %+v", courier)
	if err := pr.db.WithContext(ctx).Save(courier).Error; err != nil {
		pr.log.Errorf("Error updating courier: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("courier:%d", courier.CourierID)
	courierJSON, _ := json.Marshal(courier)
	pr.redis.Set(ctx, cacheKey, courierJSON, 0)
	pr.log.Infof("Courier saved to cache: %d", courier.CourierID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_couriers")
	pr.log.Info("Invalidated cache for all couriers")

	return courier, nil
}

func (pr *courierRepository) Delete(ctx context.Context, courierID int64) error {
	pr.log.Infof("Deleting courier with ID: %d", courierID)
	result := pr.db.WithContext(ctx).Delete(&Courier{}, courierID)
	if result.Error != nil {
		pr.log.Errorf("Error deleting courier: %v", result.Error)
		return result.Error
	}

	cacheKey := fmt.Sprintf("courier:%d", courierID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("Courier deleted from cache: %d", courierID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_couriers")
	pr.log.Info("Invalidated cache for all couriers")

	return nil
}
