package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/review/repository"
	"th3y3m/e-commerce-microservices/service/review/usecase"

	"github.com/sirupsen/logrus"
)

func NewReviewRepositoryProvider() repository.IReviewRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewReviewRepository(db, redis, log)
}

func NewReviewUsecaseProvider() usecase.IReviewUsecase {
	log := logrus.New()
	reviewRepository := NewReviewRepositoryProvider()
	return usecase.NewReviewUsecase(reviewRepository, log)
}
