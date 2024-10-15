package dependency_injection

import (
	"th3y3m/e-commerce-microservices/service/authentication/usecase"

	"github.com/sirupsen/logrus"
)

func NewAuthUsecaseProvider() usecase.IAuthUsecase {
	log := logrus.New()
	return usecase.NewAuthUsecase(log)
}
