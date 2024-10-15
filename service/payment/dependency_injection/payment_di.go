package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/payment/repository"
	"th3y3m/e-commerce-microservices/service/payment/usecase"

	"github.com/sirupsen/logrus"
)

func NewPaymentRepositoryProvider() repository.IPaymentRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewPaymentRepository(db, redis, log)
}

func NewPaymentUsecaseProvider() usecase.IPaymentUsecase {
	log := logrus.New()
	paymentRepository := NewPaymentRepositoryProvider()
	return usecase.NewPaymentUsecase(paymentRepository, log)
}
