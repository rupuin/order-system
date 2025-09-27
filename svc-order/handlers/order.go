package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"svc-order/async"
	"svc-order/dto"
	"svc-order/persistence"
)

type OrderHandler struct {
	Producer async.Producer
}

func NewOrderHandler(producer async.Producer) *OrderHandler {
	return &OrderHandler{
		Producer: producer,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq dto.Order

	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		log.Printf("Error decoding order: %v", err)
		http.Error(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}

	repo := persistence.NewRepository()
	order, err := repo.CreateOrder(orderReq.ItemID, orderReq.BuyerAddress, orderReq.Status)

	if err != nil {
		log.Printf("failed to persist order: %v", err)
		http.Error(w, fmt.Sprintf("failed to persist order: %v", err), http.StatusInternalServerError)
		return
	}

	orderEvent := order.NewEvent("order_created")

	log.Printf("publishing event: %+v", orderEvent)

	if err := h.Producer.PublishEvent(orderEvent.Key, orderEvent.Headers, orderEvent.Payload); err != nil {
		log.Printf("error publishing order event: %v", err)
		http.Error(w, fmt.Sprintf("Failed to process order: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
