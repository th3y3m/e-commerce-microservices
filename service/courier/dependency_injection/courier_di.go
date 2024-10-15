package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/courier/repository"
	"th3y3m/e-commerce-microservices/service/courier/usecase"

	"github.com/sirupsen/logrus"
)

func NewCourierRepositoryProvider() repository.ICourierRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewCourierRepository(db, redis, log)
}

func NewCourierUsecaseProvider() usecase.ICourierUsecase {
	log := logrus.New()
	courierRepository := NewCourierRepositoryProvider()
	return usecase.NewCourierUsecase(courierRepository, log)
}
