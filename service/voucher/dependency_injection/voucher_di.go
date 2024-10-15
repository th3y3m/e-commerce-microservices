package dependency_injection

import (
	"th3y3m/e-commerce-microservices/pkg/postgresql"
	redis_client "th3y3m/e-commerce-microservices/pkg/redis"
	"th3y3m/e-commerce-microservices/service/voucher/repository"
	"th3y3m/e-commerce-microservices/service/voucher/usecase"

	"github.com/sirupsen/logrus"
)

func NewVoucherRepositoryProvider() repository.IVoucherRepository {
	log := logrus.New()
	db, err := postgresql.NewGormDB()
	if err != nil {
		log.Error(err)
	}
	redis, err := redis_client.ConnectToRedis()
	if err != nil {
		log.Error(err)
	}

	return repository.NewVoucherRepository(db, redis, log)
}

func NewVoucherUsecaseProvider() usecase.IVoucherUsecase {
	log := logrus.New()
	voucherRepository := NewVoucherRepositoryProvider()
	return usecase.NewVoucherUsecase(voucherRepository, log)
}
