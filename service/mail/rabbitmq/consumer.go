package rabbitmq

import (
	"context"
	"log"
	"strconv"
	"th3y3m/e-commerce-microservices/pkg/rabbitmq"
	"th3y3m/e-commerce-microservices/service/mail/usecase"
)

func ConsumeMailNotification(ctx context.Context, mailUsecase usecase.IMailUsecase) error {
	return rabbitmq.ConsumeMessages(ctx, "order_notification_queue", func(message map[string]string) error {
		orderIdStr := message["orderId"]
		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			log.Printf("Failed to convert orderId to int64: %v", err)
			return err
		}

		url := message["url"]

		if err := mailUsecase.SendNotification(ctx, orderId, url); err != nil {
			log.Printf("Failed to send notification: %v", err)
			return err
		}

		return nil
	})
}
