package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/category/repository"
	"th3y3m/e-commerce-microservices/service/category/usecase"

	"github.com/sirupsen/logrus"
)

func NewCategoryRepositoryProvider() repository.ICategoryRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewCategoryRepository(db, redis, log)
}

func NewCategoryUsecaseProvider() usecase.ICategoryUsecase {
	log := logrus.New()
	categoryRepository := NewCategoryRepositoryProvider()
	return usecase.NewCategoryUsecase(categoryRepository, log)
}
