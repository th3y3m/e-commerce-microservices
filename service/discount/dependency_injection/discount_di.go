package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/discount/repository"
	"th3y3m/e-commerce-microservices/service/discount/usecase"

	"github.com/sirupsen/logrus"
)

func NewDiscountRepositoryProvider() repository.IDiscountRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewDiscountRepository(db, redis, log)
}

func NewDiscountUsecaseProvider() usecase.IDiscountUsecase {
	log := logrus.New()
	discountRepository := NewDiscountRepositoryProvider()
	return usecase.NewDiscountUsecase(discountRepository, log)
}
