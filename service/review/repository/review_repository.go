package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/review/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type reviewRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IReviewRepository interface {
	Get(ctx context.Context, reviewID int64) (*Review, error)
	GetAll(ctx context.Context) ([]*Review, error)
	Create(ctx context.Context, review *Review) (*Review, error)
	Update(ctx context.Context, review *Review) (*Review, error)
	Delete(ctx context.Context, reviewID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetReviewsRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetReviewsRequest) ([]*Review, error)
}

func NewReviewRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IReviewRepository {
	return &reviewRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *reviewRepository) Get(ctx context.Context, reviewID int64) (*Review, error) {
	pr.log.Infof("Fetching review with ID: %d", reviewID)
	var review Review
	cacheKey := fmt.Sprintf("review:%d", reviewID)

	// Try to get the review from Redis cache
	if pr.redis != nil {
		cachedReview, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedReview), &review); err == nil {
				pr.log.Infof("Review found in cache: %d", reviewID)
				return &review, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get review from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&review, reviewID).Error; err != nil {
		pr.log.Errorf("Error fetching review from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		reviewJSON, _ := json.Marshal(review)
		if err := pr.redis.Set(ctx, cacheKey, reviewJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save review to Redis: %v", err)
		} else {
			pr.log.Infof("Review saved to cache: %d", reviewID)
		}
	}

	return &review, nil
}

func (pr *reviewRepository) GetAll(ctx context.Context) ([]*Review, error) {
	pr.log.Info("Fetching all reviews")
	var reviews []*Review
	cacheKey := "all_reviews"

	// Try to get the reviews from Redis cache
	if pr.redis != nil {
		cachedReviews, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedReviews), &reviews); err == nil {
				pr.log.Info("Reviews found in cache")
				return reviews, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get reviews from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&reviews).Error; err != nil {
		pr.log.Errorf("Error fetching reviews from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		reviewsJSON, _ := json.Marshal(reviews)
		if err := pr.redis.Set(ctx, cacheKey, reviewsJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save reviews to Redis: %v", err)
		} else {
			pr.log.Info("Reviews saved to cache")
		}
	}

	return reviews, nil
}

func (pr *reviewRepository) Create(ctx context.Context, review *Review) (*Review, error) {
	pr.log.Infof("Creating review: %+v", review)
	if err := pr.db.WithContext(ctx).Create(review).Error; err != nil {
		pr.log.Errorf("Error creating review: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("review:%d", review.ReviewID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		reviewJSON, _ := json.Marshal(review)
		if err := pr.redis.Set(ctx, cacheKey, reviewJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save review to Redis: %v", err)
		} else {
			pr.log.Infof("Review saved to cache: %d", review.ReviewID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_reviews").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all reviews cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all reviews")
		}
	}

	// Return the newly created review (with any updated fields)
	return review, nil
}

func (pr *reviewRepository) Update(ctx context.Context, review *Review) (*Review, error) {
	pr.log.Infof("Updating review: %+v", review)
	if err := pr.db.WithContext(ctx).Save(review).Error; err != nil {
		pr.log.Errorf("Error updating review: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("review:%d", review.ReviewID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		reviewJSON, _ := json.Marshal(review)
		if err := pr.redis.Set(ctx, cacheKey, reviewJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save review to Redis: %v", err)
		} else {
			pr.log.Infof("Review saved to cache: %d", review.ReviewID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_reviews").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all reviews cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all reviews")
		}
	}

	// Return the updated review
	return review, nil
}

func (pr *reviewRepository) Delete(ctx context.Context, reviewID int64) error {
	pr.log.Infof("Deleting review with ID: %d", reviewID)
	if err := pr.db.WithContext(ctx).Delete(&Review{}, reviewID).Error; err != nil {
		pr.log.Errorf("Error deleting review: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("review:%d", reviewID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete review from Redis: %v", err)
		} else {
			pr.log.Infof("Review deleted from cache: %d", reviewID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_reviews").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all reviews cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all reviews")
		}
	}

	return nil
}

func (pr *reviewRepository) getQuerySearch(db *gorm.DB, req *model.GetReviewsRequest) *gorm.DB {
	pr.log.Infof("Building query for review search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", req.IsDeleted)
	}

	if req.ProductID != nil {
		db = db.Where("product_id = ?", req.ProductID)
	}

	if req.UserID != nil {
		db = db.Where("user_id = ?", req.UserID)
	}

	if req.MinRating != nil {
		db = db.Where("rating >= ?", req.MinRating)
	}

	if req.MaxRating != nil {
		db = db.Where("rating <= ?", req.MaxRating)
	}

	if req.Comment != "" {
		db = db.Where("comment LIKE ?", fmt.Sprintf("%%%s%%", req.Comment))
	}

	return db
}

func (pr *reviewRepository) GetList(ctx context.Context, req *model.GetReviewsRequest) ([]*Review, error) {
	pr.log.Infof("Fetching review list with request: %+v", req)
	var reviews []*Review

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

	result := db.Table("reviews").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&reviews)
	if result.Error != nil {
		pr.log.Errorf("Error fetching review list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d reviews", len(reviews))
	return reviews, nil
}
