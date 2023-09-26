package dto

type ItemTransferItemResponse struct {
	ID             int64   `json:"id,omitempty"`
	LnitmSeq       int64   `json:"lnitmseq,omitempty"`
	ItemTransferID string  `json:"item_transfer_id,omitempty"`
	ItemID         string  `json:"item_id,omitempty"`
	ItemName       string  `json:"item_name,omitempty"`
	Uom            string  `json:"uom,omitempty"`
	FullfilledQty  float64 `json:"fullfilled_qty"`
	TransferQty    float64 `json:"transfer_qty"`
	ShippedQty     float64 `json:"shipped_qty"`
	ReceiveQty     float64 `json:"receive_qty,omitempty"`
	RequestQty     float64 `json:"request_qty,omitempty"`
	ReceiveNote    string  `json:"receive_note,omitempty"`
	UnitCost       float64 `json:"unit_cost,omitempty"`
	Subtotal       float64 `json:"subtotal,omitempty"`
	Weight         float64 `json:"weight,omitempty"`
	Note           string  `json:"note,omitempty"`

	Item         *ItemResponse         `json:"item,omitempty"`
	ItemTransfer *ItemTransferResponse `json:"transfer_item,omitempty"`
}

type InTransitTransferDetailResponse struct {
	ID        string `json:"id,omitempty"`
	Code      string `json:"code,omitempty"`
	Status    int8   `json:"status,omitempty"`
	StatusStr string `json:"status_str,omitempty"`
}

type CreateItemTransferItemRequest struct {
	ItemID      int64   `json:"item_id"`
	TransferQty float64 `json:"transfer_qty"`
	RequestQty  float64 `json:"request_qty"`
	ReceiveQty  float64 `json:"receive_qty"`
	ReceiveNote string  `json:"receive_note"`
	UnitCost    float64 `json:"unit_cost"`
	Note        string  `json:"note"`
}

type UpdateItemTransferItemRequest struct {
	Id          int64   `json:"id"`
	ItemID      int64   `json:"item_id"`
	TransferQty float64 `json:"transfer_qty"`
	RequestQty  float64 `json:"request_qty"`
	ReceiveQty  float64 `json:"receive_qty"`
	ReceiveNote string  `json:"receive_note"`
	UnitCost    float64 `json:"unit_cost"`
	Note        string  `json:"note"`
}

type CreateTransferRequestDetailGPRequest struct {
	Sequence     int     `json:"lnitmseq"`
	ItemID       string  `json:"item_id" valid:"required"`
	Uom          string  `json:"uom" valid:"required"`
	RequestQty   float64 `json:"request_qty" valid:"required"`
	FulfilledQty float64 `json:"fulfill_qty" valid:"required"`
}

type UpdateTransferRequestDetailGPRequest struct {
	Lnitmseq   int     `json:"lnitmseq"`
	Itemnmbr   string  `json:"item_id" valid:"required"`
	RequestQty float64 `json:"request_qty" valid:"required"`
}

type UpdateInTransitTransferDetailGPRequest struct {
	Lnitmseq   int     `json:"lnitmseq"`
	Itemnmbr   string  `json:"item_id" valid:"required"`
	ReasonCode string  `json:"reason_code"`
	FulfillQty float64 `json:"fulfilled_qty" valid:"required"`
}

type CommitTransferRequestDetailGPRequest struct {
	Lnitmseq     int     `json:"lnitmseq"`
	FulfilledQty float64 `json:"fulfilled_qty" valid:"required"`
}
