package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/freight_rate/repository"
	"th3y3m/e-commerce-microservices/service/freight_rate/usecase"

	"github.com/sirupsen/logrus"
)

func NewFreightRateRepositoryProvider() repository.IFreightRateRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewFreightRateRepository(db, redis, log)
}

func NewFreightRateUsecaseProvider() usecase.IFreightRateUsecase {
	log := logrus.New()
	freightRateRepository := NewFreightRateRepositoryProvider()
	return usecase.NewFreightRateUsecase(freightRateRepository, log)
}
