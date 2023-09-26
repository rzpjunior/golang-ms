package dto

type ReceivingItemResponse struct {
	ID                  int64   `json:"id"`
	PurchaseOrderItemID int64   `json:"purchase_order_item_id"`
	ItemTransferItemID  int64   `json:"item_transfer_item_id"`
	DeliverQty          float64 `json:"delivery_qty"`
	RejectQty           float64 `json:"reject_qty"`
	ReceiveQty          float64 `json:"receive_qty"`
	Weight              float64 `json:"weight"`
	Note                string  `json:"note"`
	RejectReason        int8    `json:"reject_reason"`
	IsDisabled          int8    `json:"is_disabled"`

	Receiving         *ReceivingResponse         `json:"receiving"`
	PurchaseOrderItem *PurchaseOrderItemResponse `json:"purchase_order_item"`
	ItemTransferItem  *ItemTransferItemResponse  `json:"item_transfer_item"`
	Item              *ItemResponse              `json:"item"`
}

type CreateReceivingItemRequest struct {
	ReceivingItemId string  `json:"receiving_item_id"`
	ItemId          string  `json:"item_id"`
	InboundItemId   string  `json:"inbound_item_id"`
	DeliveryQty     float64 `json:"delivery_qty"`
	RejectQty       float64 `json:"reject_qty"`
	RejectReason    int8    `json:"reject_reason"`
	ReceiveQty      float64 `json:"-"`
	Note            string  `json:"note"`
}
