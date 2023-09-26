package dto

import "time"

type ReceivingResponse struct {
	ID                  int64     `json:"-"`
	SiteId              int64     `json:"site_id"`
	PurchaseOrderId     int64     `json:"purchase_order_id"`
	ItemTransferId      int64     `json:"item_transfer_id"`
	Code                string    `json:"code"`
	AtaDate             time.Time `json:"ata_date"`
	AtaTime             string    `json:"ata_time"`
	TotalWeight         float64   `json:"total_weight"`
	Note                string    `json:"note"`
	Status              int8      `json:"status"`
	InboundType         int8      `json:"inbound_type"`
	ValidSupplierReturn int8      `json:"valid_supplier_return"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           int64     `json:"created_by"`
	ConfirmedAt         time.Time `json:"confirmed_at"`
	ConfirmedBy         int64     `json:"confirmed_by"`
	Locked              int8      `json:"locked"`
	StockType           int8      `json:"stock_type"`
	LockedBy            int64     `json:"-"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           int64     `json:"-"`

	Site          *SiteResponse          `json:"site"`
	PurchaseOrder *PurchaseOrderResponse `json:"purchase_order"`
	ItemTransfer  *ItemTransferResponse  `json:"item_transfer"`

	ReceivingItems []*ReceivingItemResponse `json:"receiving_items,omitempty"`
}

type CreateReceivingRequest struct {
	Id          int64  `json:"id" valid:"required"`
	SiteId      int64  `json:"site_id" valid:"required"`
	AtaDateStr  string `json:"ata_date"`
	AtaTimeStr  string `json:"ata_time"`
	Note        string `json:"note"`
	InboundType string `json:"inbound_type" valid:"required"`

	ReceivingItem []*CreateReceivingItemRequest `json:"receiving_items"`
}

type ConfirmReceivingRequest struct {
	InboundType string `json:"inbound_type" valid:"required"`
}

type ReceivingListinDetailResponse struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Status int8   `json:"status"`
}
