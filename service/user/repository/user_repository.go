package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"th3y3m/e-commerce-microservices/service/user/model"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	log   *logrus.Logger
	db    *gorm.DB
	redis *redis.Client
}

type IUserRepository interface {
	Get(ctx context.Context, userID *int64, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, userID int64) error
	getQuerySearch(db *gorm.DB, req *model.GetUsersRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetUsersRequest) ([]*User, error)
	VerifyToken(ctx context.Context, token string, userID int64) (bool, error)
}

func NewUserRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger) IUserRepository {
	return &userRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (pr *userRepository) Get(ctx context.Context, userID *int64, email string) (*User, error) {
	var user User
	var cacheKey string

	if userID != nil {
		pr.log.Infof("Fetching user with ID: %d", *userID)
		cacheKey = fmt.Sprintf("user:%d", *userID)
	} else if email != "" {
		pr.log.Infof("Fetching user with email: %s", email)
		cacheKey = fmt.Sprintf("user:email:%s", email)
	} else {
		return nil, fmt.Errorf("either userID or email must be provided")
	}

	// Try to get the user from Redis cache
	cachedUser, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			pr.log.Infof("User found in cache: %s", cacheKey)
			return &user, nil
		}
	}

	// If not found in cache, get from database
	query := pr.db.WithContext(ctx)
	if userID != nil {
		query = query.First(&user, *userID)
	} else {
		query = query.Where("email = ?", email).First(&user)
	}

	if err := query.Error; err != nil {
		pr.log.Errorf("Error fetching user from database: %v", err)
		return nil, err
	}

	// Save to cache
	userJSON, _ := json.Marshal(user)
	pr.redis.Set(ctx, cacheKey, userJSON, 0)
	pr.log.Infof("User saved to cache: %s", cacheKey)

	return &user, nil
}

func (pr *userRepository) GetAll(ctx context.Context) ([]*User, error) {
	pr.log.Info("Fetching all users")
	var users []*User
	cacheKey := "all_users"

	// Try to get the users from Redis cache
	cachedUsers, err := pr.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedUsers), &users); err == nil {
			pr.log.Info("Users found in cache")
			return users, nil
		}
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&users).Error; err != nil {
		pr.log.Errorf("Error fetching users from database: %v", err)
		return nil, err
	}

	// Save to cache
	usersJSON, _ := json.Marshal(users)
	pr.redis.Set(ctx, cacheKey, usersJSON, 0)
	pr.log.Info("Users saved to cache")

	return users, nil
}

func (pr *userRepository) Create(ctx context.Context, user *User) (*User, error) {
	pr.log.Infof("Creating user: %+v", user)
	if err := pr.db.WithContext(ctx).Create(user).Error; err != nil {
		pr.log.Errorf("Error creating user: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("user:%d", user.UserID)

	userJSON, _ := json.Marshal(user)
	pr.redis.Set(ctx, cacheKey, userJSON, 0)
	pr.log.Infof("User saved to cache: %d", user.UserID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_users")
	pr.log.Info("Invalidated cache for all users")

	// Return the newly created user (with any updated fields)
	return user, nil
}

func (pr *userRepository) Update(ctx context.Context, user *User) (*User, error) {
	pr.log.Infof("Updating user: %+v", user)
	if err := pr.db.WithContext(ctx).Save(user).Error; err != nil {
		pr.log.Errorf("Error updating user: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("user:%d", user.UserID)

	userJSON, _ := json.Marshal(user)
	pr.redis.Set(ctx, cacheKey, userJSON, 0)
	pr.log.Infof("User saved to cache: %d", user.UserID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_users")
	pr.log.Info("Invalidated cache for all users")

	// Return the updated user
	return user, nil
}

func (pr *userRepository) Delete(ctx context.Context, userID int64) error {
	pr.log.Infof("Deleting user with ID: %d", userID)
	if err := pr.db.WithContext(ctx).Delete(&User{}, userID).Error; err != nil {
		pr.log.Errorf("Error deleting user: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("user:%d", userID)
	pr.redis.Del(ctx, cacheKey)
	pr.log.Infof("User deleted from cache: %d", userID)

	// Invalidate the cache for all records
	pr.redis.Del(ctx, "all_users")
	pr.log.Info("Invalidated cache for all users")

	return nil
}

func (pr *userRepository) getQuerySearch(db *gorm.DB, req *model.GetUsersRequest) *gorm.DB {
	pr.log.Infof("Building query for user search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", req.IsDeleted)
	}

	if req.Email != "" {
		db = db.Where("email LIKE ?", fmt.Sprintf("%%%s%%", req.Email))
	}

	if req.FullName != "" {
		db = db.Where("full_name LIKE ?", fmt.Sprintf("%%%s%%", req.FullName))
	}

	if req.PhoneNumber != "" {
		db = db.Where("phone_number LIKE ?", fmt.Sprintf("%%%s%%", req.PhoneNumber))
	}

	if req.Address != "" {
		db = db.Where("address LIKE ?", fmt.Sprintf("%%%s%%", req.Address))
	}

	if req.Role != "" {
		db = db.Where("role = ?", req.Role)
	}

	if !req.FromDate.IsZero() {
		db = db.Where("created_at >= ?", req.FromDate)
	}

	if !req.ToDate.IsZero() {
		db = db.Where("created_at <= ?", req.ToDate)
	}

	if req.Provider != "" {
		db = db.Where("provider = ?", req.Provider)
	}

	return db
}

func (pr *userRepository) GetList(ctx context.Context, req *model.GetUsersRequest) ([]*User, error) {
	pr.log.Infof("Fetching user list with request: %+v", req)
	var users []*User

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

	result := db.Table("users").Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).Limit(int(req.Paging.PageSize)).Find(&users)
	if result.Error != nil {
		pr.log.Errorf("Error fetching user list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d users", len(users))
	return users, nil
}

func (pr *userRepository) VerifyToken(ctx context.Context, token string, userID int64) (bool, error) {
	pr.log.Infof("Verifying token for user: %d", userID)

	var user User
	if err := pr.db.WithContext(ctx).Where("user_id = ? AND token = ?", userID, token).First(&user).Error; err != nil {
		pr.log.Errorf("Error fetching user with matching token from database: %v", err)
		return false, fmt.Errorf("invalid user or token")
	}

	// Check if the token has expired
	if user.TokenExpires.Before(time.Now()) {
		pr.log.Errorf("Token has expired for user: %d", userID)
		return false, fmt.Errorf("token has expired")
	}

	pr.log.Infof("Token successfully verified for user: %d", userID)
	return true, nil
}
