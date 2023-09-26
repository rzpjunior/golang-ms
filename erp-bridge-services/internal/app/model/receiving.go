package model

import "time"

type Receiving struct {
	ID                  int64     `json:"-"`
	Code                string    `json:"code"`
	SiteId              int64     `json:"site_id"`
	PurchaseOrderId     int64     `json:"purchase_order_id"`
	ItemTransferId      int64     `json:"item_transfer_id"`
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
}
