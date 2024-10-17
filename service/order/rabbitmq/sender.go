package rabbitmq

import (
	"encoding/json"
	"strconv"

	"github.com/streadway/amqp"
)

func PublishInventoryUpdateEvent(userId, cartId int64) error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
		"inventory_update_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Create the message
	message := map[string]string{
		"userId": strconv.FormatInt(userId, 10),
		"cartId": strconv.FormatInt(cartId, 10),
	}
	body, _ := json.Marshal(message)

	// Publish the message to RabbitMQ
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func PublishOrderNotificationEvent(orderId int64, url string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
		"order_notification_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	message := map[string]string{
		"orderId": strconv.FormatInt(orderId, 10),
		"url":     url,
	}
	body, _ := json.Marshal(message)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
