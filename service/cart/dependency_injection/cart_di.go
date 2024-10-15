package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/cart/repository"
	"th3y3m/e-commerce-microservices/service/cart/usecase"

	"github.com/sirupsen/logrus"
)

func NewCartRepositoryProvider() repository.ICartRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewCartRepository(db, redis, log)
}

func NewCartUsecaseProvider() usecase.ICartUsecase {
	log := logrus.New()
	cartRepository := NewCartRepositoryProvider()
	return usecase.NewCartUsecase(cartRepository, log)
}
