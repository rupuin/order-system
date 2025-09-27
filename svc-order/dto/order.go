package dto

import (
	"strconv"
	"time"
)

type Order struct {
	ID           int       `json:"id"`
	BuyerAddress string    `json:"address"`
	ItemID       string    `json:"item_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Event struct {
	Key     string            `json:"key"`
	Payload map[string]string `json:"payload"`
	Headers map[string]string `json:"headers"`
}

func (o *Order) NewEvent(eventType string) *Event {
	return &Event{
		Key: strconv.Itoa(o.ID),
		Payload: map[string]string{
			"buyer_address": o.BuyerAddress,
			"item_id":       o.ItemID,
			"status":        o.Status,
			"created_at":    o.CreatedAt.String(),
			"updated_at":    o.UpdatedAt.String(),
		},
		Headers: map[string]string{
			"event_type": eventType,
		},
	}
}
