package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cartItemRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type ICartItemRepository interface {
	Get(ctx context.Context, cartID, productID int64) (*CartItem, error)
	Create(ctx context.Context, cartItem *CartItem) (*CartItem, error)
	Update(ctx context.Context, cartItem *CartItem) (*CartItem, error)
	Delete(ctx context.Context, cartID, productID int64) error
	GetList(ctx context.Context, cartID, productID *int64) ([]*CartItem, error)
	UpdateOrCreate(ctx context.Context, cartItem CartItem) error
}

func NewCartItemRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) ICartItemRepository {
	return &cartItemRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *cartItemRepository) Get(ctx context.Context, cartID, productID int64) (*CartItem, error) {
	pr.log.Infof("Fetching cartItem with cartID: %d and productID: %d", cartID, productID)
	cacheKey := fmt.Sprintf("cartItem:%d:%d", cartID, productID)

	// Try to get the cartItem from Redis cache
	if pr.redis != nil {
		cachedCartItem, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cartItem CartItem
			if err := json.Unmarshal([]byte(cachedCartItem), &cartItem); err == nil {
				pr.log.Infof("CartItem found in cache: %+v", cartItem)
				return &cartItem, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get cartItem from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var cartItem CartItem
	result := pr.db.WithContext(ctx).First(&cartItem, "cart_id = ? AND product_id = ?", cartID, productID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching cartItem: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache with expiration if Redis is available
	if pr.redis != nil {
		cartItemJSON, _ := json.Marshal(cartItem)
		if err := pr.redis.Set(ctx, cacheKey, cartItemJSON, time.Hour).Err(); err != nil {
			pr.log.Warnf("Failed to save cartItem to Redis: %v", err)
		} else {
			pr.log.Infof("CartItem saved to cache: %d:%d", cartID, productID)
		}
	}

	return &cartItem, nil
}

func (pr *cartItemRepository) Create(ctx context.Context, cartItem *CartItem) (*CartItem, error) {
	pr.log.Infof("Creating cartItem: %+v", cartItem)
	result := pr.db.WithContext(ctx).Create(cartItem)
	if result.Error != nil {
		pr.log.Errorf("Error creating cartItem: %v", result.Error)
		return nil, result.Error
	}
	cacheKey := fmt.Sprintf("cartItem:%d:%d", cartItem.CartID, cartItem.ProductID)

	// Save to cache with expiration if Redis is available
	if pr.redis != nil {
		cartItemJSON, _ := json.Marshal(cartItem)
		if err := pr.redis.Set(ctx, cacheKey, cartItemJSON, time.Hour).Err(); err != nil {
			pr.log.Warnf("Failed to save cartItem to Redis: %v", err)
		} else {
			pr.log.Infof("CartItem saved to cache: %d:%d", cartItem.CartID, cartItem.ProductID)
		}
	}

	return cartItem, nil
}

func (pr *cartItemRepository) Update(ctx context.Context, cartItem *CartItem) (*CartItem, error) {
	pr.log.Infof("Updating cartItem: %+v", cartItem)
	result := pr.db.WithContext(ctx).Save(cartItem)
	if result.Error != nil {
		pr.log.Errorf("Error updating cartItem: %v", result.Error)
		return nil, result.Error
	}
	cacheKey := fmt.Sprintf("cartItem:%d:%d", cartItem.CartID, cartItem.ProductID)

	// Save to cache with expiration if Redis is available
	if pr.redis != nil {
		cartItemJSON, _ := json.Marshal(cartItem)
		if err := pr.redis.Set(ctx, cacheKey, cartItemJSON, time.Hour).Err(); err != nil {
			pr.log.Warnf("Failed to save cartItem to Redis: %v", err)
		} else {
			pr.log.Infof("CartItem saved to cache: %d:%d", cartItem.CartID, cartItem.ProductID)
		}
	}

	return cartItem, nil
}

func (pr *cartItemRepository) Delete(ctx context.Context, cartID, productID int64) error {
	pr.log.Infof("Deleting cartItem with cartID: %d and productID: %d", cartID, productID)
	cacheKey := fmt.Sprintf("cartItem:%d:%d", cartID, productID)

	// Delete the cartItem from the database
	result := pr.db.WithContext(ctx).Delete(&CartItem{}, "cart_id = ? AND product_id = ?", cartID, productID)
	if result.Error != nil {
		pr.log.Errorf("Error deleting cartItem: %v", result.Error)
		return result.Error
	}

	// Delete the cartItem from the cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete cartItem from Redis: %v", err)
		} else {
			pr.log.Infof("CartItem deleted from cache: %d:%d", cartID, productID)
		}
	}

	return nil
}

func (pr *cartItemRepository) GetList(ctx context.Context, cartID, productID *int64) ([]*CartItem, error) {
	pr.log.Info("Fetching cartItems")
	var cartItems []*CartItem
	cacheKey := "cartItems"

	// Build the query
	query := pr.db.WithContext(ctx)
	if cartID != nil {
		cacheKey = fmt.Sprintf("%s:cartID:%d", cacheKey, *cartID)
		query = query.Where("cart_id = ?", *cartID)
	}
	if productID != nil {
		cacheKey = fmt.Sprintf("%s:productID:%d", cacheKey, *productID)
		query = query.Where("product_id = ?", *productID)
	}

	// Try to get the cartItems from Redis cache
	if pr.redis != nil {
		cachedCartItems, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedCartItems), &cartItems); err == nil {
				pr.log.Info("CartItems found in cache")
				return cartItems, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get cartItems from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	result := query.Find(&cartItems)
	if result.Error != nil {
		pr.log.Errorf("Error fetching cartItems: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache with expiration if Redis is available
	if pr.redis != nil {
		cartItemsJSON, _ := json.Marshal(cartItems)
		if err := pr.redis.Set(ctx, cacheKey, cartItemsJSON, time.Hour).Err(); err != nil {
			pr.log.Warnf("Failed to save cartItems to Redis: %v", err)
		} else {
			pr.log.Info("CartItems saved to cache")
		}
	}

	return cartItems, nil
}

func (pr *cartItemRepository) UpdateOrCreate(ctx context.Context, cartItem CartItem) error {
	pr.log.Infof("Updating or creating cart item with cart ID: %d and product ID: %d", cartItem.CartID, cartItem.ProductID)

	var existingCartItem CartItem
	cacheKey := fmt.Sprintf("cartItem:%d:%d", cartItem.CartID, cartItem.ProductID)

	// Check if the cart item exists in the database
	if err := pr.db.WithContext(ctx).Where("cart_id = ? AND product_id = ?", cartItem.CartID, cartItem.ProductID).First(&existingCartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new cart item if not found
			pr.log.Infof("Cart item not found, creating new cart item with cart ID: %d and product ID: %d", cartItem.CartID, cartItem.ProductID)
			if err := pr.db.WithContext(ctx).Create(&cartItem).Error; err != nil {
				pr.log.Error("Failed to create new cart item:", err)
				return err
			}
		} else {
			pr.log.Error("Failed to fetch cart item for update or create:", err)
			return err
		}
	} else {
		// Update existing cart item
		pr.log.Infof("Cart item found, updating cart item with cart ID: %d and product ID: %d", cartItem.CartID, cartItem.ProductID)
		existingCartItem.Quantity = cartItem.Quantity
		if err := pr.db.WithContext(ctx).Save(&existingCartItem).Error; err != nil {
			pr.log.Error("Failed to update cart item:", err)
			return err
		}
	}

	// Save to cache with expiration if Redis is available
	if pr.redis != nil {
		cartItemJSON, _ := json.Marshal(cartItem)
		if err := pr.redis.Set(ctx, cacheKey, cartItemJSON, time.Hour).Err(); err != nil {
			pr.log.Warnf("Failed to save cartItem to Redis: %v", err)
		} else {
			pr.log.Infof("CartItem saved to cache: %d:%d", cartItem.CartID, cartItem.ProductID)
		}

		// Invalidate related cache entries
		if err := pr.redis.Del(ctx, fmt.Sprintf("cartItems:cartID:%d", cartItem.CartID)).Err(); err != nil {
			pr.log.Warnf("Failed to invalidate cartItems cache by cartID: %v", err)
		}
		if err := pr.redis.Del(ctx, fmt.Sprintf("cartItems:productID:%d", cartItem.ProductID)).Err(); err != nil {
			pr.log.Warnf("Failed to invalidate cartItems cache by productID: %v", err)
		}
	}

	pr.log.Infof("Successfully updated or created cart item with cart ID: %d and product ID: %d", cartItem.CartID, cartItem.ProductID)
	return nil
}
