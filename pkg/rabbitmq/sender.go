package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func PublishEvent(queueName string, message map[string]string) error {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	log.Println("RABBITMQ_URI: ", viper.GetString("RABBITMQ_URI"))
	conn, err := amqp.Dial(viper.GetString("RABBITMQ_URI"))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
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
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
