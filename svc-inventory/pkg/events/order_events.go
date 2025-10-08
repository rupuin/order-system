package events

type Order struct {
	ID           string `json:"id"`
	BuyerAddress string `json:"buyer_address"`
	ItemID       string `json:"item_id"`
	Status       string `json:"status"`
}

type OrderEvent struct {
	Type         string `json:"type"`
	OrderID      string `json:"order_id"`
	BuyerAddress string `json:"buyer_address"`
	ItemID       string `json:"item_id"`
	Status       string `json:"status"`
}

type Event struct {
	Key     string            `json:"key"`
	Payload map[string]string `json:"payload"`
	Headers map[string]string `json:"headers"`
}
