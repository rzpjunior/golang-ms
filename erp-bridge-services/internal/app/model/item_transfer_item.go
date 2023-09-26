package model

type ItemTransferItem struct {
	ID             int64   `json:"id"`
	ItemTransferID int64   `json:"item_transfer_id"`
	ItemID         int64   `json:"item_id"`
	DeliverQty     float64 `json:"delivery_qty"`
	ReceiveQty     float64 `json:"receive_qty"`
	RequestQty     float64 `json:"request_qty"`
	ReceiveNote    string  `json:"receive_note"`
	UnitCost       float64 `json:"unit_cost"`
	Subtotal       float64 `json:"subtotal"`
	Weight         float64 `json:"weight"`
	Note           string  `json:"note"`

	ItemTransfer *ItemTransfer `json:"transfer_item,omitempty"`
}

// TODO: init to db/gp
