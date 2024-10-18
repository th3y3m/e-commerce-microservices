package delivery

import (
	"context"
	"log"
	"th3y3m/e-commerce-microservices/service/mail/dependency_injection"
	"th3y3m/e-commerce-microservices/service/mail/rabbitmq"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers() *gin.Engine {
	r := gin.Default()

	module := dependency_injection.NewMailUsecaseProvider()
	ctx := context.Background()

	go func() {
		if err := rabbitmq.ConsumeMailNotification(ctx, module); err != nil {
			log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	mail := r.Group("/api/mail")
	{
		mail.POST("/send-mail", SendMail)
		mail.POST("/send-noti", SendNotification)
		// mail.POST("/send-order-details", SendOrderDetails)
	}

	return r
}
