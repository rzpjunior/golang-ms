package dto

import "time"

type ItemTransferResponse struct {
	ID                 int64     `json:"id"`
	Code               string    `json:"code"`
	RequestDate        time.Time `json:"request_date"`
	RecognitionDate    time.Time `json:"recognition_date"`
	EtaDate            time.Time `json:"eta_date"`
	EtaTime            string    `json:"eta_time"`
	AtaDate            time.Time `json:"ata_date"`
	AtaTime            string    `json:"ata_time"`
	AdditionalCost     float64   `json:"additional_cost"`
	AdditionalCostNote string    `json:"additional_cost_note"`
	StockType          int8      `json:"stock_type"`
	TotalCost          float64   `json:"total_cost"`
	TotalCharge        float64   `json:"total_charge"`
	TotalWeight        float64   `json:"total_weight"`
	Note               string    `json:"note"`
	Status             int8      `json:"status"`
	Locked             int8      `json:"locked"`
	LockedBy           int64     `json:"locked_by"`
	TotalSku           int64     `json:"total_sku"`
	SiteOriginID       int64     `json:"site_origin_id"`
	SiteDestinationID  int64     `json:"site_destination_id"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          int64     `json:"updated_by"`
	HasFinishedGr      int8      `json:"has_finish_gr"`

	SiteOrigin       *SiteResponse                    `json:"site_origin"`
	SiteDestination  *SiteResponse                    `json:"site_destination"`
	ItemTransferItem []*ItemTransferItemResponse      `json:"item_transfer_item,omitempty"`
	Receiving        []*ReceivingListinDetailResponse `json:"receiving"`
}

type CancelItemTransferRequest struct {
	Id   int64  `json:"-"`
	Note string `json:"note" valid:"requred"`
}

type CreateItemTransferRequest struct {
	RequestDateStr    string `json:"request_date" valid:"required"`
	SiteOriginID      int64  `json:"site_origin_id" valid:"required"`
	SiteDestinationID int64  `json:"site_destination_id" valid:"required"`
	Note              string `json:"note"`
	StockTypeID       int8   `json:"stock_type"`

	ItemTransferItems []*CreateItemTransferItemRequest `json:"item_transfer_items" valid:"required"`
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

type UpdateItemTransferRequest struct {
	Id                 int64   `json:"-"`
	SiteOriginID       int64   `json:"site_origin_id" valid:"required"`
	SiteDestinationID  int64   `json:"site_destination_id" valid:"required"`
	RecognitionDateStr string  `json:"recognition_date" valid:"required"`
	RequestDateStr     string  `json:"request_date" valid:"required"`
	EtaDateStr         string  `json:"eta_date"`
	EtaTimeStr         string  `json:"eta_time"`
	AdditionalCost     float64 `json:"additional_cost"`
	AdditionalCostNote string  `json:"additional_cost_note"`
	Note               string  `json:"note"`

	ItemTransferItems []*UpdateItemTransferItemRequest `json:"item_transfer_items" valid:"required"`
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

type CommitItemTransferRequest struct {
	Id                 int64   `json:"-"`
	RecognitionDateStr string  `json:"recognition_date" valid:"required"`
	EtaDateStr         string  `json:"eta_date" valid:"required"`
	EtaTimeStr         string  `json:"eta_time" valid:"required"`
	AdditionalCost     float64 `json:"additional_cost"`
	AdditionalCostNote string  `json:"additional_cost_note"`

	ItemTransferItems []*UpdateItemTransferItemRequest `json:"item_transfer_items" valid:"required"`
}

type CreateTransferRequestGPRequest struct {
	Interid         string                                 `json:"interid"`
	Docnumbr        string                                 `json:"docnumbr"`
	Docdate         string                                 `json:"docdate" valid:"required"`
	IvmTrType       int                                    `json:"ivm_tr_type" valid:"required"`
	RequestDate     string                                 `json:"request_date" valid:"required"`
	IvmReqEta       string                                 `json:"ivm_req_eta" valid:"required"`
	IvmLocncodeFrom string                                 `json:"ivm_locncode_from" valid:"required"`
	IvmLocncodeTo   string                                 `json:"ivm_locncode_to" valid:"required"`
	ReasonCode      string                                 `json:"reason_code"`
	Detail          []CreateTransferRequestDetailGPRequest `json:"detail" valid:"required"`
}

type UpdateTransferRequestGPRequest struct {
	Interid     string                                 `json:"interid"`
	Docnumbr    string                                 `json:"docnumbr" valid:"required"`
	Docdate     string                                 `json:"docdate" valid:"required"`
	RequestDate string                                 `json:"request_date" valid:"required"`
	IvmReqEta   string                                 `json:"ivm_req_eta" valid:"required"`
	ReasonCode  string                                 `json:"reason_code"`
	Detail      []UpdateTransferRequestDetailGPRequest `json:"detail" valid:"required"`
}

type UpdateInTransiteTransferGPRequest struct {
	Interid     string                                   `json:"interid"`
	Orddocid    string                                   `json:"orddocid" valid:"required"`
	IvmTrNumber string                                   `json:"ivm_tr_number" valid:"required"`
	Ordrdate    string                                   `json:"ordrdate" valid:"required"`
	Etadte      string                                   `json:"etadte" valid:"required"`
	Etatime     string                                   `json:"etatime" valid:"required"`
	Note        string                                   `json:"note"`
	Detail      []UpdateInTransitTransferDetailGPRequest `json:"detail" valid:"required"`
}

type CommitTransferRequestGPRequest struct {
	Interid  string                                 `json:"interid"`
	Docnumbr string                                 `json:"docnumbr"`
	ItLocn   string                                 `json:"itlocn"`
	Detail   []CommitTransferRequestDetailGPRequest `json:"detail"`
}
