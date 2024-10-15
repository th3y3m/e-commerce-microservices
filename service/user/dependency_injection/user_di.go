package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/user/repository"
	"th3y3m/e-commerce-microservices/service/user/usecase"

	"github.com/sirupsen/logrus"
)

func NewUserRepositoryProvider() repository.IUserRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewUserRepository(db, redis, log)
}

func NewUserUsecaseProvider() usecase.IUserUsecase {
	log := logrus.New()
	userRepository := NewUserRepositoryProvider()
	return usecase.NewUserUsecase(userRepository, log)
}
