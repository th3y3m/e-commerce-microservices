package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cartRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type ICartRepository interface {
	Get(ctx context.Context, cartID *int64) (*Cart, error)
	Create(ctx context.Context, cart *Cart) (*Cart, error)
	Update(ctx context.Context, cart *Cart) (*Cart, error)
	Delete(ctx context.Context, cartID int64) error
	GetUserCart(ctx context.Context, userID int64) (*Cart, error)
}

func NewCartRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) ICartRepository {
	return &cartRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *cartRepository) Get(ctx context.Context, cartID *int64) (*Cart, error) {
	pr.log.Infof("Fetching cart with ID: %d", *cartID)
	cacheKey := fmt.Sprintf("cart:%d", *cartID)

	// Try to get the cart from Redis cache
	if pr.redis != nil {
		cachedCart, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cart Cart
			if err := json.Unmarshal([]byte(cachedCart), &cart); err == nil {
				pr.log.Infof("Cart found in cache: %+v", cart)
				return &cart, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get cart from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var cart Cart
	result := pr.db.WithContext(ctx).First(&cart, cartID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching cart: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		cartJSON, _ := json.Marshal(cart)
		if err := pr.redis.Set(ctx, cacheKey, cartJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save cart to Redis: %v", err)
		} else {
			pr.log.Infof("Cart saved to cache: %d", *cartID)
		}
	}

	return &cart, nil
}

func (pr *cartRepository) Create(ctx context.Context, cart *Cart) (*Cart, error) {
	pr.log.Infof("Creating cart: %+v", cart)
	if err := pr.db.WithContext(ctx).Create(cart).Error; err != nil {
		pr.log.Errorf("Error creating cart: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("cart:%d", cart.CartID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		cartJSON, _ := json.Marshal(cart)
		if err := pr.redis.Set(ctx, cacheKey, cartJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save cart to Redis: %v", err)
		} else {
			pr.log.Infof("Cart saved to cache: %d", cart.CartID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_carts").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all carts cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all carts")
		}
	}

	return cart, nil
}

func (pr *cartRepository) Update(ctx context.Context, cart *Cart) (*Cart, error) {
	pr.log.Infof("Updating cart with ID: %d", cart.CartID)
	if err := pr.db.WithContext(ctx).Save(cart).Error; err != nil {
		pr.log.Errorf("Error updating cart: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("cart:%d", cart.CartID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		cartJSON, _ := json.Marshal(cart)
		if err := pr.redis.Set(ctx, cacheKey, cartJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save cart to Redis: %v", err)
		} else {
			pr.log.Infof("Cart saved to cache: %d", cart.CartID)
		}
	}

	return cart, nil
}

// func (pr *cartRepository) Delete(ctx context.Context, cartID int64) error {
// 	pr.log.Infof("Deleting cart with ID: %d", cartID)
// 	if err := pr.db.WithContext(ctx).Delete(&Cart{}, cartID).Error; err != nil {
// 		pr.log.Errorf("Error deleting cart: %v", err)
// 		return err
// 	}

// 	cacheKey := fmt.Sprintf("cart:%d", cartID)
// 	pr.redis.Del(ctx, cacheKey)
// 	pr.log.Infof("Cart deleted from cache: %d", cartID)

// 	// Invalidate the cache for all records
// 	pr.redis.Del(ctx, "all_carts")
// 	pr.log.Info("Invalidated cache for all carts")

//		return nil
//	}

func (pr *cartRepository) Delete(ctx context.Context, cartID int64) error {
	pr.log.Infof("Deleting cart with ID: %d", cartID)

	cart, err := pr.Get(ctx, &cartID)
	if err != nil {
		pr.log.Errorf("Error fetching cart: %v", err)
		return err
	}

	if err := pr.db.WithContext(ctx).Delete(&Cart{}, cartID).Error; err != nil {
		pr.log.Errorf("Error deleting cart: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("cart:%d", cartID)
	cacheKeyUser := fmt.Sprintf("user_cart:%d", cart.UserID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete cart from Redis: %v", err)
		} else {
			pr.log.Infof("Cart deleted from cache: %d", cartID)
		}

		if err := pr.redis.Del(ctx, cacheKeyUser).Err(); err != nil {
			pr.log.Warnf("Failed to delete user cart from Redis: %v", err)
		} else {
			pr.log.Infof("User cart deleted from cache: %d", cart.UserID)
		}
	}

	return nil
}

func (pr *cartRepository) GetUserCart(ctx context.Context, userID int64) (*Cart, error) {
	pr.log.Infof("Fetching cart for user: %d", userID)
	cacheKey := fmt.Sprintf("user_cart:%d", userID)

	// Try to get the cart from Redis cache
	if pr.redis != nil {
		cachedCart, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cart Cart
			if err := json.Unmarshal([]byte(cachedCart), &cart); err == nil {
				pr.log.Infof("Cart found in cache for user: %d", userID)
				return &cart, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get cart from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var cart Cart
	result := pr.db.WithContext(ctx).Where("user_id = ? AND is_deleted = ?", userID, false).Order("created_at DESC").First(&cart)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			pr.log.Infof("Cart not found for user: %d, creating new cart", userID)
			newCart := &Cart{
				UserID: userID,
			}
			createdCart, err := pr.Create(ctx, newCart)
			if err != nil {
				return nil, err
			}
			return createdCart, nil
		}
		pr.log.Errorf("Error fetching cart: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		cartJSON, _ := json.Marshal(cart)
		if err := pr.redis.Set(ctx, cacheKey, cartJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save cart to Redis: %v", err)
		} else {
			pr.log.Infof("Cart saved to cache for user: %d", userID)
		}
	}

	return &cart, nil
}
