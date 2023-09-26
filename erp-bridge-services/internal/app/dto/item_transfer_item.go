package dto

type ItemTransferItemResponse struct {
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

	Item         *ItemResponse         `json:"item"`
	ItemTransfer *ItemTransferResponse `json:"transfer_item,omitempty"`
}

type CreateTransferRequestDetailGPRequest struct {
	Lnitmseq      int     `json:"lnitmseq"`
	Itemnmbr      string  `json:"itemnmbr" valid:"required"`
	Uofm          string  `json:"uofm" valid:"required"`
	IvmQtyRequest float64 `json:"ivm_qty_request" valid:"required"`
	IvmQtyFulfill float64 `json:"ivm_qty_fulfill" valid:"required"`
}

type UpdateTransferRequestDetailGPRequest struct {
	Lnitmseq      int     `json:"lnitmseq"`
	Itemnmbr      string  `json:"itemnmbr" valid:"required"`
	IvmQtyRequest float64 `json:"ivm_qty_request" valid:"required"`
}

type UpdateInTransitTransferDetailGPRequest struct {
	Lnitmseq   int     `json:"lnitmseq"`
	Itemnmbr   string  `json:"itemnmbr" valid:"required"`
	ReasonCode string  `json:"reason_code" valid:"required"`
	Qtyfulfi   float64 `json:"qtyfulfi" valid:"required"`
}

type CommitTransferRequestDetailGPRequest struct {
	Lnitmseq      int     `json:"lnitmseq"`
	IvmQtyFulfill float64 `json:"ivm_qty_fulfill" valid:"required"`
}
