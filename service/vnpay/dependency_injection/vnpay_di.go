package dependency_injection

import (
	"th3y3m/e-commerce-microservices/service/vnpay/usecase"

	"github.com/sirupsen/logrus"
)

func NewVnpayUsecaseProvider() usecase.IVnpayUsecase {
	log := logrus.New()
	return usecase.NewVnpayUsecase(log)
}
