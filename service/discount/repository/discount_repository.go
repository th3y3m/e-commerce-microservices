package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type discountRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IDiscountRepository interface {
	Get(ctx context.Context, discountID int64) (*Discount, error)
	GetAll(ctx context.Context) ([]*Discount, error)
	Create(ctx context.Context, discount *Discount) (*Discount, error)
	Update(ctx context.Context, discount *Discount) (*Discount, error)
	Delete(ctx context.Context, discountID int64) error
}

func NewDiscountRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IDiscountRepository {
	return &discountRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *discountRepository) Get(ctx context.Context, discountID int64) (*Discount, error) {
	pr.log.Infof("Fetching discount with ID: %d", discountID)
	cacheKey := fmt.Sprintf("discount:%d", discountID)

	// Try to get the discount from Redis cache
	if pr.redis != nil {
		cachedDiscount, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var discount Discount
			if err := json.Unmarshal([]byte(cachedDiscount), &discount); err == nil {
				pr.log.Infof("Discount found in cache: %+v", discount)
				return &discount, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get discount from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var discount Discount
	result := pr.db.WithContext(ctx).First(&discount, discountID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching discount: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		discountJSON, _ := json.Marshal(discount)
		if err := pr.redis.Set(ctx, cacheKey, discountJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save discount to Redis: %v", err)
		} else {
			pr.log.Infof("Discount saved to cache: %d", discountID)
		}
	}

	return &discount, nil
}

func (pr *discountRepository) GetAll(ctx context.Context) ([]*Discount, error) {
	pr.log.Info("Fetching all discounts")
	var discounts []*Discount
	cacheKey := "all_discounts"

	// Try to get the discounts from Redis cache
	if pr.redis != nil {
		cachedDiscounts, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedDiscounts), &discounts); err == nil {
				pr.log.Info("Discounts found in cache")
				return discounts, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get discounts from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&discounts).Error; err != nil {
		pr.log.Errorf("Error fetching discounts from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		discountsJSON, _ := json.Marshal(discounts)
		if err := pr.redis.Set(ctx, cacheKey, discountsJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save discounts to Redis: %v", err)
		} else {
			pr.log.Info("Discounts saved to cache")
		}
	}

	return discounts, nil
}

func (pr *discountRepository) Create(ctx context.Context, discount *Discount) (*Discount, error) {
	pr.log.Infof("Creating discount: %+v", discount)
	if err := pr.db.WithContext(ctx).Create(discount).Error; err != nil {
		pr.log.Errorf("Error creating discount: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("discount:%d", discount.DiscountID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		discountJSON, _ := json.Marshal(discount)
		if err := pr.redis.Set(ctx, cacheKey, discountJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save discount to Redis: %v", err)
		} else {
			pr.log.Infof("Discount saved to cache: %d", discount.DiscountID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_discounts").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all discounts cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all discounts")
		}
	}

	// Return the newly created discount (with any updated fields)
	return discount, nil
}

func (pr *discountRepository) Update(ctx context.Context, discount *Discount) (*Discount, error) {
	pr.log.Infof("Updating discount: %+v", discount)
	if err := pr.db.WithContext(ctx).Save(discount).Error; err != nil {
		pr.log.Errorf("Error updating discount: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("discount:%d", discount.DiscountID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		discountJSON, _ := json.Marshal(discount)
		if err := pr.redis.Set(ctx, cacheKey, discountJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save discount to Redis: %v", err)
		} else {
			pr.log.Infof("Discount saved to cache: %d", discount.DiscountID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_discounts").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all discounts cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all discounts")
		}
	}

	// Return the updated discount
	return discount, nil
}

func (pr *discountRepository) Delete(ctx context.Context, discountID int64) error {
	pr.log.Infof("Deleting discount with ID: %d", discountID)
	if err := pr.db.WithContext(ctx).Delete(&Discount{}, discountID).Error; err != nil {
		pr.log.Errorf("Error deleting discount: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("discount:%d", discountID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete discount from Redis: %v", err)
		} else {
			pr.log.Infof("Discount deleted from cache: %d", discountID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_discounts").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all discounts cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all discounts")
		}
	}

	return nil
}
