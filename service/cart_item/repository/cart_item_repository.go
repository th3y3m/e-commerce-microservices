package repository

import (
	"context"
	"encoding/json"
	"fmt"

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
	cacheKey := fmt.Sprintf("cartItem:%d - %d", cartID, productID)

	// Try to get the cartItem from Redis cache
	cachedCartItem, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cartItem CartItem
		if err := json.Unmarshal([]byte(cachedCartItem), &cartItem); err == nil {
			pr.log.Infof("CartItem found in cache: %+v", cartItem)
			return &cartItem, nil
		}
	}

	// If not found in cache, get from database
	var cartItem CartItem
	result := pr.db.WithContext(ctx).First(&cartItem, "cart_id = ? AND product_id = ?", cartID, productID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching cartItem: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	cartItemJSON, _ := json.Marshal(cartItem)
	pr.redis.Set(ctx, cacheKey, cartItemJSON, 0)
	pr.log.Infof("CartItem saved to cache: %d - %d", cartID, productID)

	return &cartItem, nil
}

func (pr *cartItemRepository) Create(ctx context.Context, cartItem *CartItem) (*CartItem, error) {
	pr.log.Infof("Creating cartItem: %+v", cartItem)
	result := pr.db.WithContext(ctx).Create(cartItem)
	if result.Error != nil {
		pr.log.Errorf("Error creating cartItem: %v", result.Error)
		return nil, result.Error
	}
	cacheKey := fmt.Sprintf("cartItem:%d - %d", cartItem.CartID, cartItem.ProductID)

	// Save to cache
	cartItemJSON, _ := json.Marshal(cartItem)
	pr.redis.Set(ctx, cacheKey, cartItemJSON, 0)
	pr.log.Infof("CartItem saved to cache: %d - %d", cartItem.CartID, cartItem.ProductID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_cartItems")
	pr.log.Info("Invalidated cache for all cartItems")

	return cartItem, nil
}

func (pr *cartItemRepository) Update(ctx context.Context, cartItem *CartItem) (*CartItem, error) {
	pr.log.Infof("Updating cartItem: %+v", cartItem)
	result := pr.db.WithContext(ctx).Save(cartItem)
	if result.Error != nil {
		pr.log.Errorf("Error updating cartItem: %v", result.Error)
		return nil, result.Error
	}
	cacheKey := fmt.Sprintf("cartItem:%d - %d", cartItem.CartID, cartItem.ProductID)

	// Save to cache
	cartItemJSON, _ := json.Marshal(cartItem)
	pr.redis.Set(ctx, cacheKey, cartItemJSON, 0)
	pr.log.Infof("CartItem saved to cache: %d - %d", cartItem.CartID, cartItem.ProductID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_cartItems")
	pr.log.Info("Invalidated cache for all cartItems")

	return cartItem, nil
}

func (pr *cartItemRepository) Delete(ctx context.Context, cartID, productID int64) error {
	pr.log.Infof("Deleting cartItem with cartID: %d and productID: %d", cartID, productID)
	cacheKey := fmt.Sprintf("cartItem:%d - %d", cartID, productID)

	// Delete the cartItem from the database
	result := pr.db.WithContext(ctx).Delete(&CartItem{}, "cart_id = ? AND product_id = ?", cartID, productID)
	if result.Error != nil {
		pr.log.Errorf("Error deleting cartItem: %v", result.Error)
		return result.Error
	}

	// Delete the cartItem from the cache
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("CartItem deleted from cache: %d", cartID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_cartItems")
	pr.log.Info("Invalidated cache for all cartItems")

	return nil
}

func (pr *cartItemRepository) GetList(ctx context.Context, cartID, productID *int64) ([]*CartItem, error) {
	pr.log.Info("Fetching all cartItems")
	var cartItems []*CartItem

	// Try to get the cartItems from Redis cache
	cacheKey := "all_cartItems"
	cachedCartItems, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedCartItems), &cartItems); err == nil {
			pr.log.Info("CartItems found in cache")
			return cartItems, nil
		}
	}

	// If not found in cache, get from database
	query := pr.db.WithContext(ctx)
	if cartID != nil {
		query = query.Where("cart_id = ?", *cartID)
	}
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}
	result := query.Find(&cartItems)
	if result.Error != nil {
		pr.log.Errorf("Error fetching cartItems: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	cartItemsJSON, _ := json.Marshal(cartItems)
	pr.redis.Set(ctx, cacheKey, cartItemsJSON, 0)
	pr.log.Info("CartItems saved to cache")

	return cartItems, nil
}
