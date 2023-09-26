package dto

import "time"

type ItemTransferResponse struct {
	ID                 string    `json:"id,omitempty"`
	Code               string    `json:"code,omitempty"`
	RequestDate        time.Time `json:"request_date,omitempty"`
	RecognitionDate    time.Time `json:"recognition_date,omitempty"`
	EtaDate            time.Time `json:"eta_date,omitempty"`
	EtaTime            string    `json:"eta_time,omitempty"`
	VendorID           string    `json:"vendor_id,omitempty"`
	AtaDate            time.Time `json:"ata_date,omitempty"`
	AtaTime            string    `json:"ata_time,omitempty"`
	AdditionalCost     float64   `json:"additional_cost,omitempty"`
	AdditionalCostNote string    `json:"additional_cost_note,omitempty"`
	StockType          int8      `json:"stock_type,omitempty"`
	TotalCost          float64   `json:"total_cost,omitempty"`
	TotalCharge        float64   `json:"total_charge,omitempty"`
	TotalWeight        float64   `json:"total_weight,omitempty"`
	Note               string    `json:"note,omitempty"`
	StatusGP           int32     `json:"status_gp,omitempty"`
	StatusGPStr        string    `json:"status_gp_str,omitempty"`
	Status             int32     `json:"status,omitempty"`
	StatusStr          string    `json:"status_str,omitempty"`
	HasFinishedGr      int8      `json:"has_finished_gr,omitempty"`
	Locked             int8      `json:"locked,omitempty"`
	LockedBy           int64     `json:"locked_by,omitempty"`
	TotalSku           int64     `json:"total_sku,omitempty"`
	SiteOriginID       int64     `json:"site_origin_id,omitempty"`
	SiteDestinationID  int64     `json:"site_destination_id,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
	UpdatedBy          int64     `json:"updated_by,omitempty"`
	Type               int32     `json:"type,omitempty"`
	TypeStr            string    `json:"type_str,omitempty"`
	ReasonCode         string    `json:"reason_code,omitempty"`
	ReasonCodeStr      string    `json:"reason_code_str,omitempty"`
	ShippingMethod     string    `json:"shipping_method"`

	SiteOrigin        *SiteResponse                    `json:"site_origin,omitempty"`
	SiteDestination   *SiteResponse                    `json:"site_destination,omitempty"`
	ItemTransferItem  []*ItemTransferItemResponse      `json:"item_transfer_item,,omitempty"`
	Receiving         []*ReceivingListinDetailResponse `json:"receiving,omitempty"`
	InTransitTransfer *InTransitTransferDetailResponse `json:"intransit_transfer,omitempty"`
}

type CancelItemTransferRequest struct {
	Id   int64  `json:"-"`
	Note string `json:"note" valid:"requred"`
}

type ItemTransferListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type ItemTransferDetailRequest struct {
	Id int64 `json:"id"`
}

type CreateItemTransferRequest struct {
	RequestDateStr    string `json:"request_date" valid:"required"`
	SiteOriginID      int64  `json:"site_origin_id" valid:"required"`
	SiteDestinationID int64  `json:"site_destination_id" valid:"required"`
	Note              string `json:"note"`
	StockTypeID       int8   `json:"stock_type"`

	ItemTransferItems []*CreateItemTransferItemRequest `json:"item_transfer_items" valid:"required"`
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

type CommitItemTransferRequest struct {
	Id                 int64   `json:"-"`
	RecognitionDateStr string  `json:"recognition_date" valid:"required"`
	EtaDateStr         string  `json:"eta_date" valid:"required"`
	EtaTimeStr         string  `json:"eta_time" valid:"required"`
	AdditionalCost     float64 `json:"additional_cost"`
	AdditionalCostNote string  `json:"additional_cost_note"`

	ItemTransferItems []*UpdateItemTransferItemRequest `json:"item_transfer_items" valid:"required"`
}

type GetInTransitTransferGPListRequest struct {
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
	Orddocid    string `json:"orddocid"`
	IvmTrType   string `json:"ivm_tr_type"`
	Ordrdate    string `json:"ordrdate"`
	Trnsfloc    string `json:"trnsfloc"`
	Locncode    string `json:"locncode"`
	RequestDate string `json:"request_date"`
	Etadte      string `json:"etadte"`
	Status      int32  `json:"status"`
}

type GetTransferRequestGPListRequest struct {
	Limit           int32     `json:"limit"`
	Offset          int32     `json:"offset"`
	Docnumbr        string    `json:"docnumbr"`
	IvmTrType       string    `json:"ivm_tr_type"`
	RequestDateFrom time.Time `json:"request_date_from"`
	RequestDateTo   time.Time `json:"request_date_to"`
	IvmLocncodeFrom string    `json:"ivm_locncode_from"`
	IvmLocncodeTo   string    `json:"ivm_locncode_to"`
	DocdateFrom     string    `json:"docdate_from"`
	DocdateTo       string    `json:"docdate_to"`
	IvmReqEtaFrom   string    `json:"ivm_req_eta_from"`
	IvmReqEtaTo     string    `json:"ivm_req_eta_to"`
	Status          int32     `json:"status"`
	OrderBy         string    `query:"orderby"`
}

type CreateTransferRequestGPRequest struct {
	Docnumbr        string                                 `json:"docnumbr"`
	RecognitionDate string                                 `json:"recognition_date" valid:"required"`
	Type            int                                    `json:"type" valid:"required"`
	EtaDate         string                                 `json:"eta_date" valid:"required"`
	EtaTime         string                                 `json:"eta_time" valid:"required"`
	SiteFrom        string                                 `json:"site_from" valid:"required"`
	SiteTo          string                                 `json:"site_to" valid:"required"`
	ReasonCode      string                                 `json:"reason_code"`
	Detail          []CreateTransferRequestDetailGPRequest `json:"detail" valid:"required"`
}

type UpdateTransferRequestGPRequest struct {
	Interid     string                                 `json:"interid"`
	Docnumbr    string                                 `json:"tr_number" valid:"required"`
	Docdate     string                                 `json:"recognition_date" valid:"required"`
	RequestDate string                                 `json:"eta_date" valid:"required"`
	IvmReqEta   string                                 `json:"eta_time" valid:"required"`
	ReasonCode  string                                 `json:"reason_code"`
	Detail      []UpdateTransferRequestDetailGPRequest `json:"detail" valid:"required"`
}

type UpdateInTransiteTransferGPRequest struct {
	Interid     string                                   `json:"interid"`
	Orddocid    string                                   `json:"itt_number" valid:"required"`
	IvmTrNumber string                                   `json:"tr_number" valid:"required"`
	Ordrdate    string                                   `json:"recognition_date" valid:"required"`
	Etadte      string                                   `json:"eta_date" valid:"required"`
	Etatime     string                                   `json:"eta_time" valid:"required"`
	Note        string                                   `json:"note"`
	Detail      []UpdateInTransitTransferDetailGPRequest `json:"detail" valid:"required"`
}

type CommitTransferRequestGPRequest struct {
	TRNumber string                                 `json:"tr_number"`
	Detail   []CommitTransferRequestDetailGPRequest `json:"detail"`
}
