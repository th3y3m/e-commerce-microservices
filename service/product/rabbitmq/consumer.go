package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"th3y3m/e-commerce-microservices/service/product/usecase"

	"github.com/streadway/amqp"
)

func ConsumeInventoryUpdates(ctx context.Context, productUsecase usecase.IProductUsecase) error {
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
		"inventory_update_queue",
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

				// Deserialize message and update inventory
				var message map[string]string
				if err := json.Unmarshal(d.Body, &message); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				// Extract userId and cartId from the message
				userIdStr := message["userId"]
				cartId := message["cartId"]

				// Convert userId to int64
				userId, err := strconv.ParseInt(userIdStr, 10, 64)
				if err != nil {
					log.Printf("Failed to convert userId to int64: %v", err)
					continue
				}

				// Convert cartId to int64
				cartIdInt, err := strconv.ParseInt(cartId, 10, 64)
				if err != nil {
					log.Printf("Failed to convert cartId to int64: %v", err)
					continue
				}

				// Call the UpdateInventory method from ProductUsecase
				if err := productUsecase.UpdateInventory(ctx, userId, cartIdInt); err != nil {
					log.Printf("Failed to update inventory: %v", err)
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
