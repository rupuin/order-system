package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"svc-inventory/internal/messaging"
	"svc-inventory/internal/persistence"
	"svc-inventory/pkg/events"
)

func HandleOrderCreated(msg messaging.Message) error {
	log.Printf("Processing message: %s", string(msg.Value()))

	var orderCreated events.OrderEvent
	if err := json.Unmarshal(msg.Value(), &orderCreated); err != nil {
		return err
	}

	repo := persistence.NewRepository()
	status, err := repo.FetchItemStatus(orderCreated.ItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch item status: %w", err)
	}

	if status != "available" {
		producer := messaging.NewProducer(messaging.GetBrokers())
		itemUnavailable := events.ItemUnavailableEvent{
			OrderID: orderCreated.OrderID,
			ItemID:  orderCreated.ItemID,
			Status:  status,
		}

		log.Printf("producing message: %+v", itemUnavailable)
		err := producer.PublishEvent("item.unavailabe", orderCreated.OrderID, itemUnavailable)
		if err != nil {
			log.Printf("Failed to publish unavailable event: %v", err)
		}
		return fmt.Errorf("item id %s is unavailable", orderCreated.ItemID)
	}
	return nil
}

// func HandleOrderUpdated(message kafka.Message) error {
// 	log.Printf("Processing message: %s", string(message.Value))
//
// 	var event dto.OrderUpdatedEvent
// 	if err := json.Unmarshal(message.Value, &event); err != nil {
// 		return err
// 	}
//
// 	repo := persistence.NewRepository()
// 	orderID, err := repo.UpdateOrderStatus(event.OrderID, event.StatusTo)
// 	if err != nil {
// 		return fmt.Errorf("Failed to create order: %s", err.Error())
// 	}
//
// 	log.Printf("created order id: %d", orderID)
// 	return nil
// }
