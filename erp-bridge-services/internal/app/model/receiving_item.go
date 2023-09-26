package model

type ReceivingItem struct {
	ID                  int64   `json:"-"`
	PurchaseOrderItemID int64   `json:"purchase_order_item_id"`
	ItemTransferItemID  int64   `json:"item_transfer_item_id"`
	DeliverQty          float64 `json:"delivery_qty"`
	RejectQty           float64 `json:"reject_qty"`
	ReceiveQty          float64 `json:"receive_qty"`
	Weight              float64 `json:"weight"`
	Note                string  `json:"note"`
	RejectReason        int8    `json:"reject_reason"`
}
