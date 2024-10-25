package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/order/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type orderRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IOrderRepository interface {
	Get(ctx context.Context, orderID int64) (*Order, error)
	GetAll(ctx context.Context) ([]*Order, error)
	Create(ctx context.Context, order *Order) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
	Delete(ctx context.Context, orderID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetOrdersRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetOrdersRequest) ([]*Order, error)
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IOrderRepository {
	return &orderRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *orderRepository) Get(ctx context.Context, orderID int64) (*Order, error) {
	pr.log.Infof("Fetching order with ID: %d", orderID)
	var order Order
	cacheKey := fmt.Sprintf("order:%d", orderID)

	// Try to get the order from Redis cache
	if pr.redis != nil {
		cachedOrder, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedOrder), &order); err == nil {
				pr.log.Infof("Order found in cache: %d", orderID)
				return &order, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get order from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&order, orderID).Error; err != nil {
		pr.log.Errorf("Error fetching order from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		orderJSON, _ := json.Marshal(order)
		if err := pr.redis.Set(ctx, cacheKey, orderJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save order to Redis: %v", err)
		} else {
			pr.log.Infof("Order saved to cache: %d", orderID)
		}
	}

	return &order, nil
}

func (pr *orderRepository) GetAll(ctx context.Context) ([]*Order, error) {
	pr.log.Info("Fetching all orders")
	var orders []*Order
	cacheKey := "all_orders"

	// Try to get the orders from Redis cache
	if pr.redis != nil {
		cachedOrders, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedOrders), &orders); err == nil {
				pr.log.Info("Orders found in cache")
				return orders, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get orders from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&orders).Error; err != nil {
		pr.log.Errorf("Error fetching orders from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		ordersJSON, _ := json.Marshal(orders)
		if err := pr.redis.Set(ctx, cacheKey, ordersJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save orders to Redis: %v", err)
		} else {
			pr.log.Info("Orders saved to cache")
		}
	}

	return orders, nil
}

func (pr *orderRepository) Create(ctx context.Context, order *Order) (*Order, error) {
	pr.log.Infof("Creating order: %+v", order)
	if err := pr.db.WithContext(ctx).Create(order).Error; err != nil {
		pr.log.Errorf("Error creating order: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("order:%d", order.OrderID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		orderJSON, _ := json.Marshal(order)
		if err := pr.redis.Set(ctx, cacheKey, orderJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save order to Redis: %v", err)
		} else {
			pr.log.Infof("Order saved to cache: %d", order.OrderID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_orders").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all orders cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all orders")
		}
	}

	// Return the newly created order (with any updated fields)
	return order, nil
}

func (pr *orderRepository) Update(ctx context.Context, order *Order) (*Order, error) {
	pr.log.Infof("Updating order: %+v", order)
	if err := pr.db.WithContext(ctx).Save(order).Error; err != nil {
		pr.log.Errorf("Error updating order: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("order:%d", order.OrderID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		orderJSON, _ := json.Marshal(order)
		if err := pr.redis.Set(ctx, cacheKey, orderJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save order to Redis: %v", err)
		} else {
			pr.log.Infof("Order saved to cache: %d", order.OrderID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_orders").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all orders cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all orders")
		}
	}

	// Return the updated order
	return order, nil
}

func (pr *orderRepository) Delete(ctx context.Context, orderID int64) error {
	pr.log.Infof("Deleting order with ID: %d", orderID)
	if err := pr.db.WithContext(ctx).Delete(&Order{}, orderID).Error; err != nil {
		pr.log.Errorf("Error deleting order: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("order:%d", orderID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete order from Redis: %v", err)
		} else {
			pr.log.Infof("Order deleted from cache: %d", orderID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_orders").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all orders cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all orders")
		}
	}

	return nil
}

func (pr *orderRepository) getQuerySearch(db *gorm.DB, req *model.GetOrdersRequest) *gorm.DB {
	pr.log.Infof("Building query for order search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", req.IsDeleted)
	}

	if req.CustomerID != nil {
		db = db.Where("customer_id = ?", req.CustomerID)
	}

	if req.OrderStatus != "" {
		db = db.Where("order_status = ?", req.OrderStatus)
	}

	if req.ShippingAddress != "" {
		db = db.Where("shipping_address = ?", req.ShippingAddress)
	}

	if req.CourierID != nil {
		db = db.Where("courier_id = ?", req.CourierID)
	}

	if req.VoucherID != nil {
		db = db.Where("voucher_id = ?", req.VoucherID)
	}

	if req.MinAmount != nil {
		db = db.Where("total_amount >= ?", req.MinAmount)
	}

	if req.MaxAmount != nil {
		db = db.Where("total_amount <= ?", req.MaxAmount)
	}

	if !req.FromDate.IsZero() {
		db = db.Where("order_date >= ?", req.FromDate)
	}

	if !req.ToDate.IsZero() {
		db = db.Where("order_date <= ?", req.ToDate)
	}

	return db
}

func (pr *orderRepository) GetList(ctx context.Context, req *model.GetOrdersRequest) ([]*Order, error) {
	pr.log.Infof("Fetching order list with request: %+v", req)
	var orders []*Order

	db := pr.db.WithContext(ctx)
	db = pr.getQuerySearch(db, req)

	var sort string
	var order string

	if req.Paging.Sort == "" {
		sort = "order_date"
	} else {
		sort = req.Paging.Sort
	}

	if req.Paging.SortDirection == "" {
		order = "desc"
	} else {
		order = req.Paging.SortDirection
	}

	db = db.Order(fmt.Sprintf("%s %s", sort, order))

	result := db.Table("orders").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&orders)
	if result.Error != nil {
		pr.log.Errorf("Error fetching order list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d orders", len(orders))
	return orders, nil
}
