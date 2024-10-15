package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/order_detail/repository"
	"th3y3m/e-commerce-microservices/service/order_detail/usecase"

	"github.com/sirupsen/logrus"
)

func NewOrderDetailRepositoryProvider() repository.IOrderDetailRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewOrderDetailRepository(db, redis, log)
}

func NewOrderDetailUsecaseProvider() usecase.IOrderDetailUsecase {
	log := logrus.New()
	orderDetailRepository := NewOrderDetailRepositoryProvider()
	return usecase.NewOrderDetailUsecase(orderDetailRepository, log)
}
