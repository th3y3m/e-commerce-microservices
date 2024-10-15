package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/voucher/model"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type voucherRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IVoucherRepository interface {
	Get(ctx context.Context, voucherID int64) (*Voucher, error)
	GetAll(ctx context.Context) ([]*Voucher, error)
	Create(ctx context.Context, voucher *Voucher) (*Voucher, error)
	Update(ctx context.Context, voucher *Voucher) (*Voucher, error)
	Delete(ctx context.Context, voucherID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetVouchersRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetVouchersRequest) ([]*Voucher, error)
}

func NewVoucherRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IVoucherRepository {
	return &voucherRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *voucherRepository) Get(ctx context.Context, voucherID int64) (*Voucher, error) {
	pr.log.Infof("Fetching voucher with ID: %d", voucherID)
	var voucher Voucher
	cacheKey := fmt.Sprintf("voucher:%d", voucherID)

	// Try to get the voucher from Redis cache
	cachedVoucher, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedVoucher), &voucher); err == nil {
			pr.log.Infof("Voucher found in cache: %d", voucherID)
			return &voucher, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&voucher, voucherID).Error; err != nil {
		pr.log.Errorf("Error fetching voucher from database: %v", err)
		return nil, err
	}

	// Save to cache
	voucherJSON, _ := json.Marshal(voucher)
	pr.redis.Set(ctx, cacheKey, voucherJSON, 0)
	pr.log.Infof("Voucher saved to cache: %d", voucherID)

	return &voucher, nil
}

func (pr *voucherRepository) GetAll(ctx context.Context) ([]*Voucher, error) {
	pr.log.Info("Fetching all vouchers")
	var vouchers []*Voucher
	cacheKey := "all_vouchers"

	// Try to get the vouchers from Redis cache
	cachedVouchers, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedVouchers), &vouchers); err == nil {
			pr.log.Info("Vouchers found in cache")
			return vouchers, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&vouchers).Error; err != nil {
		pr.log.Errorf("Error fetching vouchers from database: %v", err)
		return nil, err
	}

	// Save to cache
	vouchersJSON, _ := json.Marshal(vouchers)
	pr.redis.Set(ctx, cacheKey, vouchersJSON, 0)
	pr.log.Info("Vouchers saved to cache")

	return vouchers, nil
}

func (pr *voucherRepository) Create(ctx context.Context, voucher *Voucher) (*Voucher, error) {
	pr.log.Infof("Creating voucher: %+v", voucher)
	if err := pr.db.WithContext(ctx).Create(voucher).Error; err != nil {
		pr.log.Errorf("Error creating voucher: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("voucher:%d", voucher.VoucherID)

	voucherJSON, _ := json.Marshal(voucher)
	pr.redis.Set(ctx, cacheKey, voucherJSON, 0)
	pr.log.Infof("Voucher saved to cache: %d", voucher.VoucherID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_vouchers")
	pr.log.Info("Invalidated cache for all vouchers")

	// Return the newly created voucher (with any updated fields)
	return voucher, nil
}

func (pr *voucherRepository) Update(ctx context.Context, voucher *Voucher) (*Voucher, error) {
	pr.log.Infof("Updating voucher: %+v", voucher)
	if err := pr.db.WithContext(ctx).Save(voucher).Error; err != nil {
		pr.log.Errorf("Error updating voucher: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("voucher:%d", voucher.VoucherID)

	voucherJSON, _ := json.Marshal(voucher)
	pr.redis.Set(ctx, cacheKey, voucherJSON, 0)
	pr.log.Infof("Voucher saved to cache: %d", voucher.VoucherID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_vouchers")
	pr.log.Info("Invalidated cache for all vouchers")

	// Return the updated voucher
	return voucher, nil
}

func (pr *voucherRepository) Delete(ctx context.Context, voucherID int64) error {
	pr.log.Infof("Deleting voucher with ID: %d", voucherID)
	if err := pr.db.WithContext(ctx).Delete(&Voucher{}, voucherID).Error; err != nil {
		pr.log.Errorf("Error deleting voucher: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("voucher:%d", voucherID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("Voucher deleted from cache: %d", voucherID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_vouchers")
	pr.log.Info("Invalidated cache for all vouchers")

	return nil
}

func (pr *voucherRepository) getQuerySearch(db *gorm.DB, req *model.GetVouchersRequest) *gorm.DB {
	pr.log.Infof("Building query for voucher search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", req.IsDeleted)
	}

	if req.DiscountType != "" {
		db = db.Where("discount_type = ?", req.DiscountType)
	}

	if req.IsAvailable != nil && *req.IsAvailable {
		db = db.Where("start_date <= ? AND end_date >= ? AND usage_limit > usage_count", time.Now(), time.Now())
	}

	return db
}

func (pr *voucherRepository) GetList(ctx context.Context, req *model.GetVouchersRequest) ([]*Voucher, error) {
	pr.log.Infof("Fetching voucher list with request: %+v", req)
	var vouchers []*Voucher

	db := pr.db.WithContext(ctx)
	db = pr.getQuerySearch(db, req)

	var sort string
	var order string

	if req.Paging.Sort == "" {
		sort = "created_at"
	} else {
		sort = req.Paging.Sort
	}

	if req.Paging.SortDirection == "" {
		order = "desc"
	} else {
		order = req.Paging.SortDirection
	}

	db = db.Order(fmt.Sprintf("%s %s", sort, order))

	result := db.Table("vouchers").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&vouchers)
	if result.Error != nil {
		pr.log.Errorf("Error fetching voucher list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d vouchers", len(vouchers))
	return vouchers, nil
}
