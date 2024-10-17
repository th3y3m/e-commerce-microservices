package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"th3y3m/e-commerce-microservices/service/mail/usecase"

	"github.com/streadway/amqp"
)

func ConsumeMailNotification(ctx context.Context, mailUsecase usecase.IMailUsecase) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_notification_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Process messages in a go-routine
	go func() {
		for {
			select {
			case d := <-msgs:
				log.Printf("Received a message: %s", d.Body)

				// Deserialize message and process notification
				var message map[string]string
				if err := json.Unmarshal(d.Body, &message); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				orderIdStr := message["orderId"]
				orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
				if err != nil {
					log.Printf("Failed to convert orderId to int64: %v", err)
					continue
				}

				url := message["url"]

				if err := mailUsecase.SendNotification(ctx, orderId, url); err != nil {
					log.Printf("Failed to send notification: %v", err)
				}

			case <-ctx.Done():
				log.Println("Context canceled, stopping consumer.")
				return
			}
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-ctx.Done()

	return nil
}
