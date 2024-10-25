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
	if pr.redis != nil {
		cachedCourier, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var courier Courier
			if err := json.Unmarshal([]byte(cachedCourier), &courier); err == nil {
				pr.log.Infof("Courier found in cache: %+v", courier)
				return &courier, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get courier from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var courier Courier
	result := pr.db.WithContext(ctx).First(&courier, courierID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching courier: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		courierJSON, _ := json.Marshal(courier)
		if err := pr.redis.Set(ctx, cacheKey, courierJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save courier to Redis: %v", err)
		} else {
			pr.log.Infof("Courier saved to cache: %d", courierID)
		}
	}

	return &courier, nil
}

func (pr *courierRepository) GetAll(ctx context.Context) ([]*Courier, error) {
	pr.log.Info("Fetching all couriers")
	cacheKey := "all_couriers"

	// Try to get the couriers from Redis cache
	if pr.redis != nil {
		cachedCouriers, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var couriers []*Courier
			if err := json.Unmarshal([]byte(cachedCouriers), &couriers); err == nil {
				pr.log.Infof("Couriers found in cache: %+v", couriers)
				return couriers, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get couriers from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var couriers []*Courier
	result := pr.db.WithContext(ctx).Find(&couriers)
	if result.Error != nil {
		pr.log.Errorf("Error fetching couriers: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		couriersJSON, _ := json.Marshal(couriers)
		if err := pr.redis.Set(ctx, cacheKey, couriersJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save couriers to Redis: %v", err)
		} else {
			pr.log.Info("Couriers saved to cache")
		}
	}

	return couriers, nil
}

func (pr *courierRepository) Create(ctx context.Context, courier *Courier) (*Courier, error) {
	pr.log.Infof("Creating courier: %+v", courier)
	if err := pr.db.WithContext(ctx).Create(courier).Error; err != nil {
		pr.log.Errorf("Error creating courier: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("courier:%d", courier.CourierID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		courierJSON, _ := json.Marshal(courier)
		if err := pr.redis.Set(ctx, cacheKey, courierJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save courier to Redis: %v", err)
		} else {
			pr.log.Infof("Courier saved to cache: %d", courier.CourierID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_couriers").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all couriers cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all couriers")
		}
	}

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

	// Save to cache if Redis is available
	if pr.redis != nil {
		courierJSON, _ := json.Marshal(courier)
		if err := pr.redis.Set(ctx, cacheKey, courierJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save courier to Redis: %v", err)
		} else {
			pr.log.Infof("Courier saved to cache: %d", courier.CourierID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_couriers").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all couriers cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all couriers")
		}
	}

	// Return the updated courier
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

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete courier from Redis: %v", err)
		} else {
			pr.log.Infof("Courier deleted from cache: %d", courierID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_couriers").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all couriers cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all couriers")
		}
	}

	return nil
}
