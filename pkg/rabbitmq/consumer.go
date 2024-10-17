package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func ConsumeMessages(ctx context.Context, queueName string, handler func(map[string]string) error) error {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	rabbitMQURI := viper.GetString("RABBITMQ_URI")
	log.Println("RABBITMQ_URI: ", rabbitMQURI)
	conn, err := amqp.Dial(rabbitMQURI)
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
		queueName,
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

	go func() {
		for {
			select {
			case d := <-msgs:
				var message map[string]string
				if err := json.Unmarshal(d.Body, &message); err != nil {
					fmt.Printf("Failed to unmarshal message: %v\n", err)
					continue
				}

				if err := handler(message); err != nil {
					fmt.Printf("Handler error: %v\n", err)
				}

			case <-ctx.Done():
				fmt.Println("Context canceled, stopping consumer.")
				return
			}
		}
	}()

	<-ctx.Done()
	return nil
}
