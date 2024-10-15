package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type freightRateRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IFreightRateRepository interface {
	Get(ctx context.Context, freightRateID int64) (*FreightRate, error)
	GetAll(ctx context.Context) ([]*FreightRate, error)
	Create(ctx context.Context, freightRate *FreightRate) (*FreightRate, error)
	Update(ctx context.Context, freightRate *FreightRate) (*FreightRate, error)
	Delete(ctx context.Context, freightRateID int64) error
}

func NewFreightRateRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IFreightRateRepository {
	return &freightRateRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *freightRateRepository) Get(ctx context.Context, freightRateID int64) (*FreightRate, error) {
	pr.log.Infof("Fetching freightRate with ID: %d", freightRateID)
	var freightRate FreightRate
	cacheKey := fmt.Sprintf("freightRate:%d", freightRateID)

	// Try to get the freightRate from Redis cache
	cachedFreightRate, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedFreightRate), &freightRate); err == nil {
			pr.log.Infof("FreightRate found in cache: %d", freightRateID)
			return &freightRate, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&freightRate, freightRateID).Error; err != nil {
		pr.log.Errorf("Error fetching freightRate from database: %v", err)
		return nil, err
	}

	// Save to cache
	freightRateJSON, _ := json.Marshal(freightRate)
	pr.redis.Set(ctx, cacheKey, freightRateJSON, 0)
	pr.log.Infof("FreightRate saved to cache: %d", freightRateID)

	return &freightRate, nil
}

func (pr *freightRateRepository) GetAll(ctx context.Context) ([]*FreightRate, error) {
	pr.log.Info("Fetching all freightRates")
	var freightRates []*FreightRate
	cacheKey := "all_freightRates"

	// Try to get the freightRates from Redis cache
	cachedFreightRates, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedFreightRates), &freightRates); err == nil {
			pr.log.Info("FreightRates found in cache")
			return freightRates, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&freightRates).Error; err != nil {
		pr.log.Errorf("Error fetching freightRates from database: %v", err)
		return nil, err
	}

	// Save to cache
	freightRatesJSON, _ := json.Marshal(freightRates)
	pr.redis.Set(ctx, cacheKey, freightRatesJSON, 0)
	pr.log.Info("FreightRates saved to cache")

	return freightRates, nil
}

func (pr *freightRateRepository) Create(ctx context.Context, freightRate *FreightRate) (*FreightRate, error) {
	pr.log.Infof("Creating freightRate: %+v", freightRate)
	if err := pr.db.WithContext(ctx).Create(freightRate).Error; err != nil {
		pr.log.Errorf("Error creating freightRate: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("freightRate:%d", freightRate.FreightRateID)

	freightRateJSON, _ := json.Marshal(freightRate)
	pr.redis.Set(ctx, cacheKey, freightRateJSON, 0)
	pr.log.Infof("FreightRate saved to cache: %d", freightRate.FreightRateID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_freightRates")
	pr.log.Info("Invalidated cache for all freightRates")

	// Return the newly created freightRate (with any updated fields)
	return freightRate, nil
}

func (pr *freightRateRepository) Update(ctx context.Context, freightRate *FreightRate) (*FreightRate, error) {
	pr.log.Infof("Updating freightRate: %+v", freightRate)
	if err := pr.db.WithContext(ctx).Save(freightRate).Error; err != nil {
		pr.log.Errorf("Error updating freightRate: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("freightRate:%d", freightRate.FreightRateID)

	freightRateJSON, _ := json.Marshal(freightRate)
	pr.redis.Set(ctx, cacheKey, freightRateJSON, 0)
	pr.log.Infof("FreightRate saved to cache: %d", freightRate.FreightRateID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_freightRates")
	pr.log.Info("Invalidated cache for all freightRates")

	// Return the updated freightRate
	return freightRate, nil
}

func (pr *freightRateRepository) Delete(ctx context.Context, freightRateID int64) error {
	pr.log.Infof("Deleting freightRate with ID: %d", freightRateID)
	if err := pr.db.WithContext(ctx).Delete(&FreightRate{}, freightRateID).Error; err != nil {
		pr.log.Errorf("Error deleting freightRate: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("freightRate:%d", freightRateID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("FreightRate deleted from cache: %d", freightRateID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_freightRates")
	pr.log.Info("Invalidated cache for all freightRates")

	return nil
}
