package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/order/repository"
	"th3y3m/e-commerce-microservices/service/order/usecase"

	"github.com/sirupsen/logrus"
)

func NewOrderRepositoryProvider() repository.IOrderRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewOrderRepository(db, redis, log)
}

func NewOrderUsecaseProvider() usecase.IOrderUsecase {
	log := logrus.New()
	orderRepository := NewOrderRepositoryProvider()
	return usecase.NewOrderUsecase(orderRepository, log)
}
