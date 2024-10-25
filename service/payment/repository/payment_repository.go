package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/payment/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type paymentRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IPaymentRepository interface {
	Get(ctx context.Context, paymentID int64) (*Payment, error)
	GetAll(ctx context.Context) ([]*Payment, error)
	Create(ctx context.Context, payment *Payment) (*Payment, error)
	Update(ctx context.Context, payment *Payment) (*Payment, error)
	Delete(ctx context.Context, paymentID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetPaymentsRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetPaymentsRequest) ([]*Payment, error)
}

func NewPaymentRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IPaymentRepository {
	return &paymentRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *paymentRepository) Get(ctx context.Context, paymentID int64) (*Payment, error) {
	pr.log.Infof("Fetching payment with ID: %d", paymentID)
	var payment Payment
	cacheKey := fmt.Sprintf("payment:%d", paymentID)

	// Try to get the payment from Redis cache
	if pr.redis != nil {
		cachedPayment, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedPayment), &payment); err == nil {
				pr.log.Infof("Payment found in cache: %d", paymentID)
				return &payment, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get payment from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&payment, paymentID).Error; err != nil {
		pr.log.Errorf("Error fetching payment from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		paymentJSON, _ := json.Marshal(payment)
		if err := pr.redis.Set(ctx, cacheKey, paymentJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save payment to Redis: %v", err)
		} else {
			pr.log.Infof("Payment saved to cache: %d", paymentID)
		}
	}

	return &payment, nil
}

func (pr *paymentRepository) GetAll(ctx context.Context) ([]*Payment, error) {
	pr.log.Info("Fetching all payments")
	var payments []*Payment
	cacheKey := "all_payments"

	// Try to get the payments from Redis cache
	if pr.redis != nil {
		cachedPayments, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedPayments), &payments); err == nil {
				pr.log.Info("Payments found in cache")
				return payments, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get payments from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&payments).Error; err != nil {
		pr.log.Errorf("Error fetching payments from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		paymentsJSON, _ := json.Marshal(payments)
		if err := pr.redis.Set(ctx, cacheKey, paymentsJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save payments to Redis: %v", err)
		} else {
			pr.log.Info("Payments saved to cache")
		}
	}

	return payments, nil
}

func (pr *paymentRepository) Create(ctx context.Context, payment *Payment) (*Payment, error) {
	pr.log.Infof("Creating payment: %+v", payment)
	if err := pr.db.WithContext(ctx).Create(payment).Error; err != nil {
		pr.log.Errorf("Error creating payment: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("payment:%d", payment.PaymentID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		paymentJSON, _ := json.Marshal(payment)
		if err := pr.redis.Set(ctx, cacheKey, paymentJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save payment to Redis: %v", err)
		} else {
			pr.log.Infof("Payment saved to cache: %d", payment.PaymentID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_payments").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all payments cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all payments")
		}
	}

	// Return the newly created payment (with any updated fields)
	return payment, nil
}

func (pr *paymentRepository) Update(ctx context.Context, payment *Payment) (*Payment, error) {
	pr.log.Infof("Updating payment: %+v", payment)
	if err := pr.db.WithContext(ctx).Save(payment).Error; err != nil {
		pr.log.Errorf("Error updating payment: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("payment:%d", payment.PaymentID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		paymentJSON, _ := json.Marshal(payment)
		if err := pr.redis.Set(ctx, cacheKey, paymentJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save payment to Redis: %v", err)
		} else {
			pr.log.Infof("Payment saved to cache: %d", payment.PaymentID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_payments").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all payments cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all payments")
		}
	}

	// Return the updated payment
	return payment, nil
}

func (pr *paymentRepository) Delete(ctx context.Context, paymentID int64) error {
	pr.log.Infof("Deleting payment with ID: %d", paymentID)
	if err := pr.db.WithContext(ctx).Delete(&Payment{}, paymentID).Error; err != nil {
		pr.log.Errorf("Error deleting payment: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("payment:%d", paymentID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete payment from Redis: %v", err)
		} else {
			pr.log.Infof("Payment deleted from cache: %d", paymentID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_payments").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all payments cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all payments")
		}
	}

	return nil
}

func (pr *paymentRepository) getQuerySearch(db *gorm.DB, req *model.GetPaymentsRequest) *gorm.DB {
	pr.log.Infof("Building query for payment search: %+v", req)

	if req.PaymentStatus != "" {
		db = db.Where("payment_status = ?", req.PaymentStatus)
	}

	if req.PaymentMethod != "" {
		db = db.Where("payment_method = ?", req.PaymentMethod)
	}

	if !req.FromDate.IsZero() {
		db = db.Where("payment_date >= ?", req.FromDate)
	}

	if !req.ToDate.IsZero() {
		db = db.Where("payment_date <= ?", req.ToDate)
	}

	if req.MinAmount != nil {
		db = db.Where("payment_amount >= ?", req.MinAmount)
	}

	if req.MaxAmount != nil {
		db = db.Where("payment_amount <= ?", req.MaxAmount)
	}

	if req.OrderID != nil {
		db = db.Where("order_id = ?", req.OrderID)
	}

	return db
}

func (pr *paymentRepository) GetList(ctx context.Context, req *model.GetPaymentsRequest) ([]*Payment, error) {
	pr.log.Infof("Fetching payment list with request: %+v", req)
	var payments []*Payment

	db := pr.db.WithContext(ctx)
	db = pr.getQuerySearch(db, req)

	var sort string
	var order string

	if req.Paging.Sort == "" {
		sort = "payment_date"
	} else {
		sort = req.Paging.Sort
	}

	if req.Paging.SortDirection == "" {
		order = "desc"
	} else {
		order = req.Paging.SortDirection
	}

	db = db.Order(fmt.Sprintf("%s %s", sort, order))

	result := db.Table("payments").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&payments)
	if result.Error != nil {
		pr.log.Errorf("Error fetching payment list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d payments", len(payments))
	return payments, nil
}
