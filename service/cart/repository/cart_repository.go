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
	cachedCart, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cart Cart
		if err := json.Unmarshal([]byte(cachedCart), &cart); err == nil {
			pr.log.Infof("Cart found in cache: %+v", cart)
			return &cart, nil
		}
	}

	// If not found in cache, get from database
	var cart Cart
	result := pr.db.WithContext(ctx).First(&cart, cartID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching cart: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	cartJSON, _ := json.Marshal(cart)
	pr.redis.Set(ctx, cacheKey, cartJSON, 0)
	pr.log.Infof("Cart saved to cache: %d", *cartID)

	return &cart, nil
}

func (pr *cartRepository) Create(ctx context.Context, cart *Cart) (*Cart, error) {
	pr.log.Infof("Creating cart: %+v", cart)
	if err := pr.db.WithContext(ctx).Create(cart).Error; err != nil {
		pr.log.Errorf("Error creating cart: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("cart:%d", cart.CartID)
	cartJSON, _ := json.Marshal(cart)
	pr.redis.Set(ctx, cacheKey, cartJSON, 0)
	pr.log.Infof("Cart saved to cache: %d", cart.CartID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_carts")
	pr.log.Info("Invalidated cache for all carts")

	return cart, nil
}

func (pr *cartRepository) Update(ctx context.Context, cart *Cart) (*Cart, error) {
	pr.log.Infof("Updating cart with ID: %d", cart.CartID)
	if err := pr.db.WithContext(ctx).Save(cart).Error; err != nil {
		pr.log.Errorf("Error updating cart: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("cart:%d", cart.CartID)
	cartJSON, _ := json.Marshal(cart)
	pr.redis.Set(ctx, cacheKey, cartJSON, 0)
	pr.log.Infof("Cart saved to cache: %d", cart.CartID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_carts")
	pr.log.Info("Invalidated cache for all carts")

	return cart, nil
}

func (pr *cartRepository) Delete(ctx context.Context, cartID int64) error {
	pr.log.Infof("Deleting cart with ID: %d", cartID)
	if err := pr.db.WithContext(ctx).Delete(&Cart{}, cartID).Error; err != nil {
		pr.log.Errorf("Error deleting cart: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("cart:%d", cartID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("Cart deleted from cache: %d", cartID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_carts")
	pr.log.Info("Invalidated cache for all carts")

	return nil
}