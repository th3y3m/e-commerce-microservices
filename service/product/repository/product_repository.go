package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/product/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IProductRepository interface {
	Get(ctx context.Context, productID int64) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Create(ctx context.Context, product *Product) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	Delete(ctx context.Context, productID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetProductsRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetProductsRequest) ([]*Product, error)
}

func NewProductRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IProductRepository {
	return &productRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *productRepository) Get(ctx context.Context, productID int64) (*Product, error) {
	pr.log.Infof("Fetching product with ID: %d", productID)
	var product Product
	cacheKey := fmt.Sprintf("product:%d", productID)

	// Try to get the product from Redis cache
	cachedProduct, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
			pr.log.Infof("Product found in cache: %d", productID)
			return &product, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		pr.log.Errorf("Error fetching product from database: %v", err)
		return nil, err
	}

	// Save to cache
	productJSON, _ := json.Marshal(product)
	pr.redis.Set(ctx, cacheKey, productJSON, 0)
	pr.log.Infof("Product saved to cache: %d", productID)

	return &product, nil
}

func (pr *productRepository) GetAll(ctx context.Context) ([]*Product, error) {
	pr.log.Info("Fetching all products")
	var products []*Product
	cacheKey := "all_products"

	// Try to get the products from Redis cache
	cachedProducts, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedProducts), &products); err == nil {
			pr.log.Info("Products found in cache")
			return products, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&products).Error; err != nil {
		pr.log.Errorf("Error fetching products from database: %v", err)
		return nil, err
	}

	// Save to cache
	productsJSON, _ := json.Marshal(products)
	pr.redis.Set(ctx, cacheKey, productsJSON, 0)
	pr.log.Info("Products saved to cache")

	return products, nil
}

func (pr *productRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	pr.log.Infof("Creating product: %+v", product)
	if err := pr.db.WithContext(ctx).Create(product).Error; err != nil {
		pr.log.Errorf("Error creating product: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("product:%d", product.ProductID)

	productJSON, _ := json.Marshal(product)
	pr.redis.Set(ctx, cacheKey, productJSON, 0)
	pr.log.Infof("Product saved to cache: %d", product.ProductID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_products")
	pr.log.Info("Invalidated cache for all products")

	// Return the newly created product (with any updated fields)
	return product, nil
}

func (pr *productRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	pr.log.Infof("Updating product: %+v", product)
	if err := pr.db.WithContext(ctx).Save(product).Error; err != nil {
		pr.log.Errorf("Error updating product: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("product:%d", product.ProductID)

	productJSON, _ := json.Marshal(product)
	pr.redis.Set(ctx, cacheKey, productJSON, 0)
	pr.log.Infof("Product saved to cache: %d", product.ProductID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_products")
	pr.log.Info("Invalidated cache for all products")

	// Return the updated product
	return product, nil
}

func (pr *productRepository) Delete(ctx context.Context, productID int64) error {
	pr.log.Infof("Deleting product with ID: %d", productID)
	if err := pr.db.WithContext(ctx).Delete(&Product{}, productID).Error; err != nil {
		pr.log.Errorf("Error deleting product: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", productID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("Product deleted from cache: %d", productID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_products")
	pr.log.Info("Invalidated cache for all products")

	return nil
}

func (pr *productRepository) getQuerySearch(db *gorm.DB, req *model.GetProductsRequest) *gorm.DB {
	pr.log.Infof("Building query for product search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", *req.IsDeleted)
	}

	if req.SellerID != nil {
		db = db.Where("seller_id = ?", *req.SellerID)
	}

	if req.CategoryID != nil {
		db = db.Where("category_id = ?", *req.CategoryID)
	}

	if req.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+req.ProductName+"%")
	}

	if req.Description != "" {
		db = db.Where("description LIKE ?", "%"+req.Description+"%")
	}

	if req.MinPrice != nil {
		db = db.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		db = db.Where("price <= ?", *req.MaxPrice)
	}

	if req.MinQuantity != nil {
		db = db.Where("quantity >= ?", *req.MinQuantity)
	}

	if req.MaxQuantity != nil {
		db = db.Where("quantity <= ?", *req.MaxQuantity)
	}

	if req.FromDate.IsZero() == false {
		db = db.Where("created_at >= ?", req.FromDate)
	}

	if req.ToDate.IsZero() == false {
		db = db.Where("created_at <= ?", req.ToDate)
	}

	return db
}

func (pr *productRepository) GetList(ctx context.Context, req *model.GetProductsRequest) ([]*Product, error) {
	pr.log.Infof("Fetching product list with request: %+v", req)
	var products []*Product

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

	result := db.Table("products").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&products)
	if result.Error != nil {
		pr.log.Errorf("Error fetching product list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d products", len(products))
	return products, nil
}
