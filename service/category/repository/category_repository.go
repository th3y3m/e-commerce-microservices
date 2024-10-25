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
	if pr.redis != nil {
		cachedCategory, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var category Category
			if err := json.Unmarshal([]byte(cachedCategory), &category); err == nil {
				pr.log.Infof("Category found in cache: %+v", category)
				return &category, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get category from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var category Category
	result := pr.db.WithContext(ctx).First(&category, categoryID)
	if result.Error != nil {
		pr.log.Errorf("Error fetching category: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		categoryJSON, _ := json.Marshal(category)
		if err := pr.redis.Set(ctx, cacheKey, categoryJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save category to Redis: %v", err)
		} else {
			pr.log.Infof("Category saved to cache: %d", categoryID)
		}
	}

	return &category, nil
}

func (pr *categoryRepository) GetAll(ctx context.Context) ([]*Category, error) {
	pr.log.Info("Fetching all categories")
	cacheKey := "all_categories"

	// Try to get the categories from Redis cache
	if pr.redis != nil {
		cachedCategories, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var categories []*Category
			if err := json.Unmarshal([]byte(cachedCategories), &categories); err == nil {
				pr.log.Infof("Categories found in cache: %d", len(categories))
				return categories, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get categories from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	var categories []*Category
	result := pr.db.WithContext(ctx).Find(&categories)
	if result.Error != nil {
		pr.log.Errorf("Error fetching categories: %v", result.Error)
		return nil, result.Error
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		categoriesJSON, _ := json.Marshal(categories)
		if err := pr.redis.Set(ctx, cacheKey, categoriesJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save categories to Redis: %v", err)
		} else {
			pr.log.Info("Categories saved to cache")
		}
	}

	return categories, nil
}

func (pr *categoryRepository) Create(ctx context.Context, category *Category) (*Category, error) {
	pr.log.Infof("Creating category: %+v", category)
	if err := pr.db.WithContext(ctx).Create(category).Error; err != nil {
		pr.log.Errorf("Error creating category: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("category:%d", category.CategoryID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		categoryJSON, _ := json.Marshal(category)
		if err := pr.redis.Set(ctx, cacheKey, categoryJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save category to Redis: %v", err)
		} else {
			pr.log.Infof("Category saved to cache: %d", category.CategoryID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_categories").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all categories cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all categories")
		}
	}

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

	// Save to cache if Redis is available
	if pr.redis != nil {
		categoryJSON, _ := json.Marshal(category)
		if err := pr.redis.Set(ctx, cacheKey, categoryJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save category to Redis: %v", err)
		} else {
			pr.log.Infof("Category saved to cache: %d", category.CategoryID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_categories").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all categories cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all categories")
		}
	}

	// Return the updated category
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

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete category from Redis: %v", err)
		} else {
			pr.log.Infof("Deleted category with ID: %d", categoryID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_categories").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all categories cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all categories")
		}
	}

	return nil
}
