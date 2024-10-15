package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/cart_item/repository"
	"th3y3m/e-commerce-microservices/service/cart_item/usecase"

	"github.com/sirupsen/logrus"
)

func NewCartItemRepositoryProvider() repository.ICartItemRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewCartItemRepository(db, redis, log)
}

func NewCartItemUsecaseProvider() usecase.ICartItemUsecase {
	log := logrus.New()
	cartItemRepository := NewCartItemRepositoryProvider()
	return usecase.NewCartItemUsecase(cartItemRepository, log)
}
