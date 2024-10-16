package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/order_detail/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type orderDetailRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IOrderDetailRepository interface {
	Get(ctx context.Context, orderID, productID int64) (*OrderDetail, error)
	Create(ctx context.Context, orderDetail *OrderDetail) (*OrderDetail, error)
	Update(ctx context.Context, orderDetail *OrderDetail) (*OrderDetail, error)
	Delete(ctx context.Context, orderID, productID int64) error
	GetList(ctx context.Context, req *model.GetOrderDetailsRequest) ([]*OrderDetail, error)
}

func NewOrderDetailRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IOrderDetailRepository {
	return &orderDetailRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *orderDetailRepository) Get(ctx context.Context, orderID, productID int64) (*OrderDetail, error) {
	pr.log.Infof("Fetching orderDetail with ID: %d and productID: %d", orderID, productID)
	var orderDetail OrderDetail
	cacheKey := fmt.Sprintf("orderDetail:%d - %d", orderID, productID)

	// Try to get the orderDetail from Redis cache
	cachedOrderDetail, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedOrderDetail), &orderDetail); err == nil {
			pr.log.Infof("OrderDetail found in cache: %d:%d", orderID, productID)
			return &orderDetail, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).First(&orderDetail).Error; err != nil {
		pr.log.Errorf("Error fetching orderDetail from database: %v", err)
		return nil, err
	}

	// Save to cache
	orderDetailJSON, _ := json.Marshal(orderDetail)
	pr.redis.Set(ctx, cacheKey, orderDetailJSON, 0)
	pr.log.Infof("OrderDetail saved to cache: %d", orderID)

	return &orderDetail, nil
}

func (pr *orderDetailRepository) Create(ctx context.Context, orderDetail *OrderDetail) (*OrderDetail, error) {
	pr.log.Infof("Creating orderDetail: %+v", orderDetail)
	if err := pr.db.WithContext(ctx).Create(orderDetail).Error; err != nil {
		pr.log.Errorf("Error creating orderDetail: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("orderDetail:%d - %d", orderDetail.OrderID, orderDetail.ProductID)

	orderDetailJSON, _ := json.Marshal(orderDetail)
	pr.redis.Set(ctx, cacheKey, orderDetailJSON, 0)
	pr.log.Infof("OrderDetail saved to cache: %d - %d", orderDetail.OrderID, orderDetail.ProductID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_orderDetails")
	pr.log.Info("Invalidated cache for all orderDetails")

	return orderDetail, nil
}

func (pr *orderDetailRepository) Update(ctx context.Context, orderDetail *OrderDetail) (*OrderDetail, error) {
	pr.log.Infof("Updating orderDetail: %+v", orderDetail)
	if err := pr.db.WithContext(ctx).Save(orderDetail).Error; err != nil {
		pr.log.Errorf("Error updating orderDetail: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("orderDetail:%d - %d", orderDetail.OrderID, orderDetail.ProductID)

	orderDetailJSON, _ := json.Marshal(orderDetail)
	pr.redis.Set(ctx, cacheKey, orderDetailJSON, 0)
	pr.log.Infof("OrderDetail saved to cache: %d - %d", orderDetail.OrderID, orderDetail.ProductID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_orderDetails")
	pr.log.Info("Invalidated cache for all orderDetails")

	// Return the updated orderDetail
	return orderDetail, nil
}

func (pr *orderDetailRepository) Delete(ctx context.Context, orderID, productID int64) error {
	pr.log.Infof("Deleting orderDetail with ID: %d and productID: %d", orderID, productID)
	orderDetail, err := pr.Get(ctx, orderID, productID)
	if err != nil {
		pr.log.Errorf("Error fetching orderDetail for deletion: %v", err)
		return err
	}

	if err := pr.db.WithContext(ctx).Delete(orderDetail).Error; err != nil {
		pr.log.Errorf("Error deleting orderDetail: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("orderDetail:%d - %d", orderID, productID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("OrderDetail deleted from cache: %d", orderID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_orderDetails")
	pr.log.Info("Invalidated cache for all orderDetails")

	return nil
}

func (pr *orderDetailRepository) GetList(ctx context.Context, req *model.GetOrderDetailsRequest) ([]*OrderDetail, error) {
	pr.log.Infof("Fetching orderDetails with request: %+v", req)
	var orderDetails []*OrderDetail
	cacheKey := "all_orderDetails"

	// Try to get the orderDetails from Redis cache
	cachedOrderDetails, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedOrderDetails), &orderDetails); err == nil {
			pr.log.Infof("OrderDetails found in cache: %d", req.OrderID)
			return orderDetails, nil
		}
	}

	// If not found in cache, get from database
	if req.OrderID != nil && req.ProductID == nil {
		if err := pr.db.WithContext(ctx).Where("order_id = ?", req.OrderID).Find(&orderDetails).Error; err != nil {
			pr.log.Errorf("Error fetching orderDetails from database: %v", err)
			return nil, err
		}

	} else if req.OrderID == nil && req.ProductID != nil {
		if err := pr.db.WithContext(ctx).Where("product_id = ?", req.ProductID).Find(&orderDetails).Error; err != nil {
			pr.log.Errorf("Error fetching orderDetails from database: %v", err)
			return nil, err
		}
	} else {
		if err := pr.db.WithContext(ctx).Find(&orderDetails).Error; err != nil {
			pr.log.Errorf("Error fetching orderDetails from database: %v", err)
			return nil, err
		}
	}

	return orderDetails, nil
}
