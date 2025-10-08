package events

type ItemUnavailableEvent struct {
	OrderID string `json:"order_id"`
	ItemID  string `json:"item_id"`
	Status  string `json:"status"`
}
