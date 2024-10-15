package dependency_injection

import (
	"th3y3m/e-commerce-microservices/service/mail/usecase"

	"github.com/sirupsen/logrus"
)

func NewMailUsecaseProvider() usecase.IMailUsecase {
	log := logrus.New()
	return usecase.NewMailUsecase(log)
}
