package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/news/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type newRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type INewRepository interface {
	Get(ctx context.Context, newID int64) (*News, error)
	GetAll(ctx context.Context) ([]*News, error)
	Create(ctx context.Context, new *News) (*News, error)
	Update(ctx context.Context, new *News) (*News, error)
	Delete(ctx context.Context, newID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetNewsRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetNewsRequest) ([]*News, error)
}

func NewNewsRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) INewRepository {
	return &newRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *newRepository) Get(ctx context.Context, newID int64) (*News, error) {
	pr.log.Infof("Fetching new with ID: %d", newID)
	var new News
	cacheKey := fmt.Sprintf("new:%d", newID)

	// Try to get the new from Redis cache
	if pr.redis != nil {
		cachedNew, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedNew), &new); err == nil {
				pr.log.Infof("New found in cache: %d", newID)
				return &new, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get new from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&new, newID).Error; err != nil {
		pr.log.Errorf("Error fetching new from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		newJSON, _ := json.Marshal(new)
		if err := pr.redis.Set(ctx, cacheKey, newJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save new to Redis: %v", err)
		} else {
			pr.log.Infof("New saved to cache: %d", newID)
		}
	}

	return &new, nil
}

func (pr *newRepository) GetAll(ctx context.Context) ([]*News, error) {
	pr.log.Info("Fetching all news")
	var news []*News
	cacheKey := "all_news"

	// Try to get the news from Redis cache
	if pr.redis != nil {
		cachedNews, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedNews), &news); err == nil {
				pr.log.Info("News found in cache")
				return news, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get news from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&news).Error; err != nil {
		pr.log.Errorf("Error fetching news from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		newsJSON, _ := json.Marshal(news)
		if err := pr.redis.Set(ctx, cacheKey, newsJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save news to Redis: %v", err)
		} else {
			pr.log.Info("News saved to cache")
		}
	}

	return news, nil
}

func (pr *newRepository) Create(ctx context.Context, new *News) (*News, error) {
	pr.log.Infof("Creating new: %+v", new)
	if err := pr.db.WithContext(ctx).Create(new).Error; err != nil {
		pr.log.Errorf("Error creating new: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("new:%d", new.NewsID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		newJSON, _ := json.Marshal(new)
		if err := pr.redis.Set(ctx, cacheKey, newJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save new to Redis: %v", err)
		} else {
			pr.log.Infof("New saved to cache: %d", new.NewsID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_news").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all news cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all news")
		}
	}

	// Return the newly created new (with any updated fields)
	return new, nil
}

func (pr *newRepository) Update(ctx context.Context, new *News) (*News, error) {
	pr.log.Infof("Updating new: %+v", new)
	if err := pr.db.WithContext(ctx).Save(new).Error; err != nil {
		pr.log.Errorf("Error updating new: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("new:%d", new.NewsID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		newJSON, _ := json.Marshal(new)
		if err := pr.redis.Set(ctx, cacheKey, newJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save new to Redis: %v", err)
		} else {
			pr.log.Infof("New saved to cache: %d", new.NewsID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_news").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all news cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all news")
		}
	}

	// Return the updated new
	return new, nil
}

func (pr *newRepository) Delete(ctx context.Context, newID int64) error {
	pr.log.Infof("Deleting new with ID: %d", newID)
	if err := pr.db.WithContext(ctx).Delete(&News{}, newID).Error; err != nil {
		pr.log.Errorf("Error deleting new: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("new:%d", newID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete new from Redis: %v", err)
		} else {
			pr.log.Infof("New deleted from cache: %d", newID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_news").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all news cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all news")
		}
	}

	return nil
}

func (pr *newRepository) getQuerySearch(db *gorm.DB, req *model.GetNewsRequest) *gorm.DB {
	pr.log.Infof("Building query for new search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", req.IsDeleted)
	}

	if !req.FromDate.IsZero() {
		db = db.Where("created_at >= ?", req.FromDate)
	}

	if !req.ToDate.IsZero() {
		db = db.Where("created_at <= ?", req.ToDate)
	}

	return db
}

func (pr *newRepository) GetList(ctx context.Context, req *model.GetNewsRequest) ([]*News, error) {
	pr.log.Infof("Fetching new list with request: %+v", req)
	var news []*News

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

	result := db.Table("news").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&news)
	if result.Error != nil {
		pr.log.Errorf("Error fetching new list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d news", len(news))
	return news, nil
}
