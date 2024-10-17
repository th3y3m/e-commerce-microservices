package rabbitmq

import (
	"strconv"
	"th3y3m/e-commerce-microservices/pkg/rabbitmq"
)

func PublishInventoryUpdateEvent(userId, cartId int64) error {
	return rabbitmq.PublishEvent("inventory_update_queue", map[string]string{
		"userId": strconv.FormatInt(userId, 10),
		"cartId": strconv.FormatInt(cartId, 10),
	})
}

func PublishOrderNotificationEvent(orderId int64, url string) error {
	return rabbitmq.PublishEvent("order_notification_queue", map[string]string{
		"orderId": strconv.FormatInt(orderId, 10),
		"url":     url,
	})
}
