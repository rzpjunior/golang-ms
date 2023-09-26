package dto

import "time"

type ReceivingResponse struct {
	ID                  string    `json:"id,omitempty"`
	PurchaseOrderId     int64     `json:"purchase_order_id,omitempty"`
	ItemTransferId      int64     `json:"item_transfer_id,omitempty"`
	Code                string    `json:"code,omitempty"`
	AtaDate             time.Time `json:"ata_date,omitempty"`
	AtaTime             string    `json:"ata_time,omitempty"`
	TotalWeight         float64   `json:"total_weight,omitempty"`
	Note                string    `json:"note,omitempty"`
	Status              int8      `json:"status,omitempty"`
	StatusStr           string    `json:"status_str,omitempty"`
	Region              string    `json:"region,omitempty"`
	InboundType         int8      `json:"inbound_type,omitempty"`
	ValidSupplierReturn int8      `json:"valid_supplier_return,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	CreatedBy           int64     `json:"created_by,omitempty"`
	ConfirmedAt         time.Time `json:"confirmed_at,omitempty"`
	ConfirmedBy         int64     `json:"confirmed_by,omitempty"`
	Locked              int8      `json:"locked,omitempty"`
	StockType           int8      `json:"stock_type,omitempty"`
	LockedBy            int64     `json:"-"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
	UpdatedBy           int64     `json:"-"`

	Site          *SiteResponse          `json:"site"`
	Vendor        *VendorResponse        `json:"vendor"`
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

	ReceivingItem []*CreateReceivingItemRequest `json:"goods_receipt_item"`
}

type ReceivingListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type ReceivingDetailRequest struct {
	Id int64 `json:"id"`
}

type ConfirmReceivingRequest struct {
	Id          string `json:"-"`
	InboundType string `json:"inbound_type" valid:"required"`
}

type ReceivingListinDetailResponse struct {
	ID        string `json:"id,omitempty"`
	Code      string `json:"code,omitempty"`
	Status    int8   `json:"status,omitempty"`
	StatusStr string `json:"status_str,omitempty"`
}

type PurchaseInvoiceDetailResponse struct {
	ID        string `json:"id,omitempty"`
	Code      string `json:"code,omitempty"`
	Status    int8   `json:"status,omitempty"`
	StatusStr string `json:"status_str,omitempty"`
}

type GetGoodsReceiptGPListRequest struct {
	Limit    int32  `query:"limit"`
	Offset   int32  `query:"offset"`
	Poprctnm string `query:"poprctnm"`
	Doctype  string `query:"doctype"`
}

type GoodsReceiptDTL struct {
	Poprctnm string  `json:"poprctnm"`
	Ponumber string  `json:"ponumber"`
	Locncode string  `json:"locncode"`
	Uofm     string  `json:"uofm"`
	Itemnmbr string  `json:"itemnmbr"`
	Qtyshppd float32 `json:"qtyshppd"`
	Unitcost float64 `json:"unitcost"`
	Extdcost float64 `json:"extdcost"`
}

type CreateGoodsReceiptGPRequest struct {
	Interid         string            `json:"interid"`
	Poprctnm        string            `json:"poprctnm"`
	Vnddocnm        string            `json:"vnddocnm"`
	Receiptdate     string            `json:"receiptdate"`
	Vendorid        string            `json:"vendorid"`
	Curncyid        string            `json:"curncyid"`
	Subtotal        float64           `json:"subtotal"`
	GoodsReceiptDTL []GoodsReceiptDTL `json:"goodsReceiptDTL"`
}

type UpdateDetailGoodsReceiptGPRequest struct {
	Rcptlnnm int     `json:"rcptlnnm"`
	Ponumber string  `json:"ponumber"`
	Locncode string  `json:"locncode"`
	Uofm     string  `json:"uofm"`
	Itemnmbr string  `json:"itemnmbr"`
	Qtyshppd float32 `json:"qtyshppd"`
	Unitcost float64 `json:"unitcost"`
	Extdcost float64 `json:"extdcost"`
}

type UpdateGoodsReceiptGPRequest struct {
	Interid              string                              `json:"interid"`
	Poprctnm             string                              `json:"poprctnm"`
	PrpRegion            string                              `json:"prp_region"`
	Note                 string                              `json:"note"`
	Actlship             string                              `json:"actlship"`
	PrpActualarrivalTime string                              `json:"prp_actualarrival_time"`
	Details              []UpdateDetailGoodsReceiptGPRequest `json:"details"`
}
