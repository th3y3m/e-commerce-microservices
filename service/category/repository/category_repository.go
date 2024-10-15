package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type categoryRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type ICategoryRepository interface {
	Get(ctx context.Context, categoryID int64) (*Category, error)
	GetAll(ctx context.Context) ([]*Category, error)
	Create(ctx context.Context, category *Category) (*Category, error)
	Update(ctx context.Context, category *Category) (*Category, error)
	Delete(ctx context.Context, categoryID int64) error
}

func NewCategoryRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) ICategoryRepository {
	return &categoryRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *categoryRepository) Get(ctx context.Context, categoryID int64) (*Category, error) {
	pr.log.Infof("Fetching category with ID: %d", categoryID)
	cacheKey := fmt.Sprintf("category:%d", categoryID)

	// Try to get the category from Redis cache
	cachedCategory, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var category Category
		if err := json.Unmarshal([]byte(cachedCategory), &category); err == nil {
			pr.log.Infof("Category found in cache: %+v", category)
			return &category, nil
		}
	}

	// If not found in cache, get from database
	var category Category
	result := pr.db.WithContext(ctx).First(&category, categoryID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching category: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	categoryJSON, _ := json.Marshal(category)
	pr.redis.Set(ctx, cacheKey, categoryJSON, 0)
	pr.log.Infof("Category saved to cache: %d", categoryID)

	return &category, nil
}

func (pr *categoryRepository) GetAll(ctx context.Context) ([]*Category, error) {
	pr.log.Info("Fetching all categorys")
	cacheKey := "all_categorys"

	// Try to get the categorys from Redis cache
	cachedCategorys, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var categorys []*Category
		if err := json.Unmarshal([]byte(cachedCategorys), &categorys); err == nil {
			pr.log.Infof("Categorys found in cache: %d", len(categorys))
			return categorys, nil
		}
	}

	// If not found in cache, get from database
	var categorys []*Category
	result := pr.db.WithContext(ctx).Find(&categorys)
	if result.Error != nil {
		pr.log.Errorf("Error fetching categorys: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache
	categorysJSON, _ := json.Marshal(categorys)
	pr.redis.Set(ctx, cacheKey, categorysJSON, 0)
	pr.log.Info("Categorys saved to cache")

	return categorys, nil
}

func (pr *categoryRepository) Create(ctx context.Context, category *Category) (*Category, error) {
	pr.log.Infof("Creating category: %+v", category)
	if err := pr.db.WithContext(ctx).Create(category).Error; err != nil {
		pr.log.Errorf("Error creating category: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("category:%d", category.CategoryID)

	categoryJSON, _ := json.Marshal(category)
	pr.redis.Set(ctx, cacheKey, categoryJSON, 0)
	pr.log.Infof("Category saved to cache: %d", category.CategoryID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_categorys")
	pr.log.Info("Invalidated cache for all categorys")

	// Return the newly created category (with any updated fields)
	return category, nil
}

func (pr *categoryRepository) Update(ctx context.Context, category *Category) (*Category, error) {
	pr.log.Infof("Updating category: %+v", category)
	if err := pr.db.WithContext(ctx).Save(category).Error; err != nil {
		pr.log.Errorf("Error updating category: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("category:%d", category.CategoryID)
	categoryJSON, _ := json.Marshal(category)
	pr.redis.Set(ctx, cacheKey, categoryJSON, 0)
	pr.log.Infof("Category saved to cache: %d", category.CategoryID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_categorys")
	pr.log.Info("Invalidated cache for all categorys")

	return category, nil
}

func (pr *categoryRepository) Delete(ctx context.Context, categoryID int64) error {
	pr.log.Infof("Deleting category with ID: %d", categoryID)
	result := pr.db.WithContext(ctx).Delete(&Category{}, categoryID)
	if result.Error != nil {
		pr.log.Errorf("Error deleting category: %v", result.Error)
		return result.Error
	}

	cacheKey := fmt.Sprintf("category:%d", categoryID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("Deleted category with ID: %d", categoryID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_categorys")
	pr.log.Info("Invalidated cache for all categorys")

	return nil
}
