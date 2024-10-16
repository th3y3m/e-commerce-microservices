package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/product_discount/repository"
	"th3y3m/e-commerce-microservices/service/product_discount/usecase"

	"github.com/sirupsen/logrus"
)

func NewProductDiscountRepositoryProvider() repository.IProductDiscountRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewProductDiscountRepository(db, redis, log)
}

func NewProductDiscountUsecaseProvider() usecase.IProductDiscountUsecase {
	log := logrus.New()
	productDiscountRepository := NewProductDiscountRepositoryProvider()
	return usecase.NewProductDiscountUsecase(productDiscountRepository, log)
}
