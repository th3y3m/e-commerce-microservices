package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/news/repository"
	"th3y3m/e-commerce-microservices/service/news/usecase"

	"github.com/sirupsen/logrus"
)

func NewNewsRepositoryProvider() repository.INewRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewNewsRepository(db, redis, log)
}

func NewNewsUsecaseProvider() usecase.INewUsecase {
	log := logrus.New()
	newRepository := NewNewsRepositoryProvider()
	return usecase.NewNewsUsecase(newRepository, log)
}
