package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/service/product_discount/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productDiscountRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IProductDiscountRepository interface {
	Create(ctx context.Context, productDiscount *ProductDiscount) (*ProductDiscount, error)
	Delete(ctx context.Context, productID, discountID int64) error
	Get(ctx context.Context, req *model.GetProductDiscountsRequest) ([]*ProductDiscount, error)
}

func NewProductDiscountRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IProductDiscountRepository {
	return &productDiscountRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *productDiscountRepository) Create(ctx context.Context, productDiscount *ProductDiscount) (*ProductDiscount, error) {
	pr.log.Infof("Creating productDiscount: %+v", productDiscount)
	if err := pr.db.WithContext(ctx).Create(productDiscount).Error; err != nil {
		pr.log.Errorf("Error creating productDiscount: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("productDiscount:%d:%d", productDiscount.ProductID, productDiscount.DiscountID)

	productDiscountJSON, _ := json.Marshal(productDiscount)
	if pr.redis != nil {
		if err := pr.redis.Set(ctx, cacheKey, productDiscountJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save productDiscount to Redis: %v", err)
		} else {
			pr.log.Infof("ProductDiscount saved to cache: %d:%d", productDiscount.ProductID, productDiscount.DiscountID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_productDiscounts").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all productDiscounts cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all productDiscounts")
		}
	}

	return productDiscount, nil
}

func (pr *productDiscountRepository) Delete(ctx context.Context, productID, discountID int64) error {
	pr.log.Infof("Deleting productDiscount with productID: %d and discountID: %d", productID, discountID)
	var req model.GetProductDiscountsRequest = model.GetProductDiscountsRequest{
		ProductID:  &productID,
		DiscountID: &discountID,
	}
	productDiscounts, err := pr.Get(ctx, &req)
	if err != nil {
		pr.log.Errorf("Error fetching productDiscount for deletion: %v", err)
		return err
	}

	if len(productDiscounts) == 0 {
		pr.log.Warnf("No productDiscount found for productID: %d and discountID: %d", productID, discountID)
		return nil
	}

	if err := pr.db.WithContext(ctx).Delete(productDiscounts[0]).Error; err != nil {
		pr.log.Errorf("Error deleting productDiscount: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("productDiscount:%d:%d", productID, discountID)
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete productDiscount from Redis: %v", err)
		} else {
			pr.log.Infof("ProductDiscount deleted from cache: %d:%d", productID, discountID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_productDiscounts").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all productDiscounts cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all productDiscounts")
		}
	}

	return nil
}

func (pr *productDiscountRepository) Get(ctx context.Context, req *model.GetProductDiscountsRequest) ([]*ProductDiscount, error) {
	pr.log.Infof("Fetching productDiscounts with request: %+v", req)
	var productDiscounts []*ProductDiscount
	var cacheKey string

	db := pr.db.WithContext(ctx)

	if req.ProductID != nil && req.DiscountID == nil {
		pr.log.Infof("Fetching discounts for product ID: %d", *req.ProductID)

		cacheKey = fmt.Sprintf("productDiscounts:productID:%d", *req.ProductID)
		if pr.redis != nil {
			cachedData, err := pr.redis.Get(ctx, cacheKey).Result()
			if err == nil {
				if err := json.Unmarshal([]byte(cachedData), &productDiscounts); err == nil {
					pr.log.Infof("Fetched product discounts for product ID: %d from cache", *req.ProductID)
					return productDiscounts, nil
				}
			} else if err != redis.Nil {
				pr.log.Warnf("Failed to get product discounts from Redis: %v", err)
			}
		} else {
			pr.log.Warn("Redis client is not initialized")
		}
		db = db.Where("product_id = ?", *req.ProductID)
	} else if req.ProductID == nil && req.DiscountID != nil {
		pr.log.Infof("Fetching products for discount ID: %d", *req.DiscountID)

		cacheKey = fmt.Sprintf("productDiscounts:discountID:%d", *req.DiscountID)
		if pr.redis != nil {
			cachedData, err := pr.redis.Get(ctx, cacheKey).Result()
			if err == nil {
				if err := json.Unmarshal([]byte(cachedData), &productDiscounts); err == nil {
					pr.log.Infof("Fetched product discounts for discount ID: %d from cache", *req.DiscountID)
					return productDiscounts, nil
				}
			} else if err != redis.Nil {
				pr.log.Warnf("Failed to get product discounts from Redis: %v", err)
			}
		} else {
			pr.log.Warn("Redis client is not initialized")
		}
		db = db.Where("discount_id = ?", *req.DiscountID)
	} else if req.ProductID != nil && req.DiscountID != nil {
		pr.log.Infof("Fetching specific product discount for product ID: %d and discount ID: %d", *req.ProductID, *req.DiscountID)

		cacheKey = fmt.Sprintf("productDiscounts:productID:%d:discountID:%d", *req.ProductID, *req.DiscountID)
		if pr.redis != nil {
			cachedData, err := pr.redis.Get(ctx, cacheKey).Result()
			if err == nil {
				if err := json.Unmarshal([]byte(cachedData), &productDiscounts); err == nil {
					pr.log.Infof("Fetched specific product discount for product ID: %d and discount ID: %d from cache", *req.ProductID, *req.DiscountID)
					return productDiscounts, nil
				}
			} else if err != redis.Nil {
				pr.log.Warnf("Failed to get product discounts from Redis: %v", err)
			}
		} else {
			pr.log.Warn("Redis client is not initialized")
		}
		db = db.Where("product_id = ? AND discount_id = ?", *req.ProductID, *req.DiscountID)
	} else {
		pr.log.Info("Fetching all product discounts")
		cacheKey = "all_productDiscounts"
		if pr.redis != nil {
			cachedData, err := pr.redis.Get(ctx, cacheKey).Result()
			if err == nil {
				if err := json.Unmarshal([]byte(cachedData), &productDiscounts); err == nil {
					pr.log.Info("Fetched all product discounts from cache")
					return productDiscounts, nil
				}
			} else if err != redis.Nil {
				pr.log.Warnf("Failed to get product discounts from Redis: %v", err)
			}
		} else {
			pr.log.Warn("Redis client is not initialized")
		}
	}

	if err := db.Find(&productDiscounts).Error; err != nil {
		pr.log.Errorf("Error fetching product discounts: %v", err)
		return nil, err
	}

	if len(productDiscounts) == 0 {
		pr.log.Infof("No product discounts found for request: %+v", req)
		return nil, constant.ErrNoProductDiscountsFound
	}

	productDiscountsJSON, _ := json.Marshal(productDiscounts)
	if pr.redis != nil {
		if err := pr.redis.Set(ctx, cacheKey, productDiscountsJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save product discounts to Redis: %v", err)
		} else {
			pr.log.Infof("Product discounts saved to cache: %s", cacheKey)
		}
	}

	pr.log.Infof("Fetched %d product discounts", len(productDiscounts))
	return productDiscounts, nil
}
