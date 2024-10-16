package dependency_injection

import (
	"th3y3m/e-commerce-microservices/service/momo/usecase"

	"github.com/sirupsen/logrus"
)

func NewMoMoUsecaseProvider() usecase.IMoMoUsecase {
	log := logrus.New()
	return usecase.NewMoMoUsecase(log)
}
