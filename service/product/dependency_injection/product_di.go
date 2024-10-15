package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/product/repository"
	"th3y3m/e-commerce-microservices/service/product/usecase"

	"github.com/sirupsen/logrus"
)

func NewProductRepositoryProvider() repository.IProductRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewProductRepository(db, redis, log)
}

func NewProductUsecaseProvider() usecase.IProductUsecase {
	log := logrus.New()
	productRepository := NewProductRepositoryProvider()
	return usecase.NewProductUsecase(productRepository, log)
}
