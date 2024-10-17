package rabbitmq

import (
	"context"
	"log"
	"strconv"
	"th3y3m/e-commerce-microservices/pkg/rabbitmq"
	"th3y3m/e-commerce-microservices/service/product/usecase"
)

func ConsumeInventoryUpdates(ctx context.Context, productUsecase usecase.IProductUsecase) error {
	return rabbitmq.ConsumeMessages(ctx, "inventory_update_queue", func(message map[string]string) error {
		userIdStr := message["userId"]
		cartIdStr := message["cartId"]

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			log.Printf("Failed to convert userId to int64: %v", err)
			return err
		}

		cartId, err := strconv.ParseInt(cartIdStr, 10, 64)
		if err != nil {
			log.Printf("Failed to convert cartId to int64: %v", err)
			return err
		}

		if err := productUsecase.UpdateInventory(ctx, userId, cartId); err != nil {
			log.Printf("Failed to update inventory: %v", err)
			return err
		}

		return nil
	})
}
