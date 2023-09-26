package dto

type ReceivingItemResponse struct {
	ID                  string  `json:"id,omitempty"`
	ItemName            string  `json:"item_name,omitempty"`
	Uom                 string  `json:"uom,omitempty"`
	PurchaseOrderItemID int64   `json:"purchase_order_item_id,omitempty"`
	ItemTransferItemID  int64   `json:"item_transfer_item_id,omitempty"`
	DeliverQty          float64 `json:"delivery_qty"`
	RejectQty           float64 `json:"reject_qty"`
	ShippedQty          float64 `json:"shipped_qty"`
	OrderQty            float64 `json:"order_qty"`
	TransferQty         float64 `json:"transfer_qty"`
	Weight              float64 `json:"weight,omitempty"`
	Note                string  `json:"note,omitempty"`
	RejectReason        int8    `json:"reject_reason,omitempty"`
	RejectReasonConvert string  `json:"reject_reason_convert,omitempty"`
	IsDisabled          int8    `json:"is_disabled,omitempty"`

	Receiving         *ReceivingResponse         `json:"goods_receipt"`
	PurchaseOrderItem *PurchaseOrderItemResponse `json:"purchase_order_item"`
	ItemTransferItem  *ItemTransferItemResponse  `json:"item_transfer_item"`
	Item              *ItemResponse              `json:"item"`
}

type CreateReceivingItemRequest struct {
	ItemId        int64   `json:"item_id"`
	InboundItemId int64   `json:"inbound_item_id"`
	DeliveryQty   float64 `json:"delivery_qty"`
	RejectQty     float64 `json:"reject_qty"`
	RejectReason  int8    `json:"reject_reason"`
	ReceiveQty    float64 `json:"-"`
	Note          string  `json:"note"`
}
