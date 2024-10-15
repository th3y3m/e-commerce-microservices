package dependency_injection

import (
	"th3y3m/e-commerce-microservices/service/oauth/usecase"

	"github.com/sirupsen/logrus"
)

func NewOAuthUsecaseProvider() usecase.IOAuthUsecase {
	log := logrus.New()
	return usecase.NewOAuthUsecase(log)
}
