package dto

import "time"

type PurchaseOrderResponse struct {
	ID                     int64     `json:"id"`
	Code                   string    `json:"code"`
	VendorID               int64     `json:"vendor_id"`
	SiteID                 int64     `json:"site_id"`
	TermPaymentPurID       int64     `json:"term_payment_pur_id"`
	VendorClassificationID int64     `json:"vendor_classification_id"`
	PurchasePlanID         int64     `json:"purchase_plan_id"`
	ConsolidatedShipmentID int64     `json:"consolidated_shipment_id"`
	Status                 int32     `json:"status"`
	DocDate                time.Time `json:"doc_date"`
	RecognitionDate        time.Time `json:"recognition_date"`
	EtaDate                time.Time `json:"eta_date"`
	SiteAddress            string    `json:"site_address"`
	EtaTime                string    `json:"eta_time"`
	TaxPct                 float64   `json:"tax_pct"`
	DeliveryFee            float64   `json:"delivery_fee"`
	TotalPrice             float64   `json:"total_price"`
	TaxAmount              float64   `json:"tax_amount"`
	TotalCharge            float64   `json:"total_charge"`
	TotalInvoice           float64   `json:"total_invoice"`
	TotalWeight            float64   `json:"total_weight"`
	Note                   string    `json:"note"`
	DeltaPrint             int32     `json:"delta_print"`
	Latitude               float64   `json:"latitude"`
	Longitude              float64   `json:"longitude"`
	UpdatedAt              time.Time `json:"updated_at"`
	UpdatedBy              int64     `json:"updated_by"`
	CreatedFrom            int32     `json:"created_from"`
	HasFinishedGr          int8      `json:"has_finished_gr"`
	CreatedAt              time.Time `json:"created_at"`
	CreatedBy              int64     `json:"created_by"`
	CommittedAt            time.Time `json:"committed_at"`
	CommittedBy            int64     `json:"committed_by"`
	AssignedTo             int64     `json:"assigned_to"`
	AssignedBy             int64     `json:"assigned_by"`
	AssignedAt             time.Time `json:"assigned_at"`
	Locked                 int32     `json:"locked"`
	LockedBy               int64     `json:"locked_by"`

	PurchaseOrderItems []*PurchaseOrderItemResponse     `json:"purchase_order_items"`
	Receiving          []*ReceivingListinDetailResponse `json:"receiving"`
}

type CreatePurchaseOrderGPRequest struct {
	Interid                 string   `json:"interid"`
	Potype                  int64    `json:"potype"`
	Ponumber                string   `json:"ponumber"`
	Docdate                 string   `json:"docdate"`
	Buyerid                 string   `json:"buyerid"`
	Vendorid                string   `json:"vendorid"`
	Curncyid                string   `json:"curncyid"`
	Deprtmnt                string   `json:"deprtmnt"`
	Locncode                string   `json:"locncode"`
	Taxschid                string   `json:"taxschid"`
	Subtotal                float64  `json:"subtotal"`
	Trdisamt                float64  `json:"trdisamt"`
	Frtamnt                 float64  `json:"frtamnt"`
	Miscamnt                float64  `json:"miscamnt"`
	Taxamnt                 float64  `json:"taxamnt"`
	PrpPurchaseplanNo       string   `json:"prp_purchaseplan_no"`
	CsReference             struct{} `json:"cs_reference"`
	PrpPaymentMethod        string   `json:"prp_payment_method"`
	PrpPaymentTerm          string   `json:"pymtrmid"`
	DueDate                 string   `json:"duedate"`
	PrpRegion               string   `json:"prp_region"`
	PrpEstimatedarrivalDate string   `json:"prp_estimatedarrival_date"`
	Notetext                string   `json:"notetext"`
	PRStatus                int32    `json:"pr_status"`
	Detail                  []PODTL  `json:"detail"`
}

type UpdatePurchaseOrderGPRequest struct {
	Interid                 string                         `json:"interid"`
	Potype                  int64                          `json:"potype"`
	Ponumber                string                         `json:"ponumber"`
	Docdate                 string                         `json:"docdate"`
	Buyerid                 string                         `json:"buyerid"`
	Vendorid                string                         `json:"vendorid"`
	Curncyid                string                         `json:"curncyid"`
	Deprtmnt                string                         `json:"deprtmnt"`
	Locncode                string                         `json:"locncode"`
	Taxschid                string                         `json:"taxschid"`
	Subtotal                float64                        `json:"subtotal"`
	Trdisamt                float64                        `json:"trdisamt"`
	Frtamnt                 float64                        `json:"frtamnt"`
	Miscamnt                float64                        `json:"miscamnt"`
	Taxamnt                 float64                        `json:"taxamnt"`
	PrpPurchaseplanNo       string                         `json:"prp_purchaseplan_no"`
	CsReference             *ConsolidatedShipmentReference `json:"cs_reference"`
	Pymtrmid                string                         `json:"pymtrmid"`
	Duedate                 string                         `json:"duedate"`
	PrpPaymentMethod        string                         `json:"prp_payment_method"`
	PrpRegion               string                         `json:"prp_region"`
	PrpEstimatedarrivalDate string                         `json:"prp_estimatedarrival_date"`
	Notetext                string                         `json:"notetext"`
	Detail                  []PODTL                        `json:"detail"`
}

type CreatePurchaseOrderRequest struct {
	VendorID    int64   `json:"vendor_id" valid:"required"`
	SiteID      int64   `json:"site_id" valid:"required"`
	OrderDate   string  `json:"order_date" valid:"required"`
	StrEtaDate  string  `json:"eta_date" valid:"required"`
	EtaTime     string  `json:"eta_time" valid:"required"`
	DeliveryFee float64 `json:"delivery_fee"`
	Note        string  `json:"note"`
	TaxPct      float64 `json:"tax_pct"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`

	PurchaseOrderItems []*CreatePurchaseOrderItemRequest `json:"purchase_order_items" valid:"required"`
}

type PODTL struct {
	Ord      int32   `json:"ord"`
	Itemnmbr string  `json:"itemnmbr"`
	Uofm     string  `json:"uofm"`
	Qtyorder float64 `json:"qtyorder"`
	Qtycance float64 `json:"qtycance"`
	Unitcost float64 `json:"unitcost"`
	Notetext string  `json:"notetext"`
}

type UpdatePurchaseOrderRequest struct {
	Id            int64   `json:"id"`
	VendorID      int64   `json:"vendor_id" valid:"required"`
	SiteID        int64   `json:"site_id" valid:"required"`
	OrderDate     string  `json:"order_date" valid:"required"`
	StrEtaDate    string  `json:"eta_date" valid:"required"`
	EtaTime       string  `json:"eta_time" valid:"required"`
	DeliveryFee   float64 `json:"delivery_fee"`
	PaymentTermID string  `json:"payment_term_id" valid:"required"`
	Note          string  `json:"note"`
	TaxPct        float64 `json:"tax_pct"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`

	PurchaseOrderItems []*UpdatePurchaseOrderItemRequest `json:"purchase_order_items" valid:"required"`
}

type UpdateProductPurchaseOrderRequest struct {
	Id          int64   `json:"-"`
	DeliveryFee float64 `json:"delivery_fee"`
	TaxPct      float64 `json:"tax_pct"`

	PurchaseOrderItems []*UpdatePurchaseOrderItemRequest `json:"purchase_order_items" valid:"required"`
}

type CreatePurchaseOrderItemRequest struct {
	ItemID        int64   `json:"item_id" valid:"required"`
	OrderQty      float64 `json:"qty" valid:"required"`
	UnitPrice     float64 `json:"unit_price" valid:"required"`
	Note          string  `json:"note" valid:"lte:500"`
	PurchaseQty   float64 `json:"purchase_qty"`
	IncludeTax    int8    `json:"include_tax"`
	TaxPercentage float64 `json:"tax_percentage"`
}

type UpdatePurchaseOrderItemRequest struct {
	Id            int64   `json:"id" valid:"required"`
	ItemID        int64   `json:"item_id" valid:"required"`
	OrderQty      float64 `json:"qty" valid:"required"`
	UnitPrice     float64 `json:"unit_price" valid:"required"`
	Note          string  `json:"note" valid:"lte:500"`
	PurchaseQty   float64 `json:"purchase_qty"`
	IncludeTax    int8    `json:"include_tax"`
	TaxPercentage float64 `json:"tax_percentage"`
}

type CancelPurchaseOrderRequest struct {
	Id   int64  `json:"-"`
	Note string `json:"note" valid:"requred"`
}

type ConsolidatedShipmentReference struct {
	PRPCSNo      string `json:"prp_cs_no"`
	PRVehicleNo  string `json:"prp_vehicle_no"`
	PRDriverName string `json:"prp_driver_name"`
	PhonName     string `json:"phonname"`
}

type CancelPurchaseOrderGPRequest struct {
	Interid             string `json:"interid"`
	PurchaseOrderNumber string `json:"PONUMBER"`
	UserId              string `json:"userid"`
}

type CancelPurchaseOrderGPResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NEW CONSOLIDATED SHIPMENT
type POLST struct {
	Ponumber string `json:"ponumber"`
}

type CreateConsolidatedShipmentGPRequest struct {
	Interid         string `json:"interid"`
	PRPCSNo         string `json:"prp_cs_no"`
	PrpDriverName   string `json:"prp_driver_name" valid:"required"`
	PrVehicleNumber string `json:"pr_vehicle_number" valid:"required"`
	PhoneName       string `json:"phonname" valid:"required"`

	PurchaseOrders []POLST `json:"polist" valid:"required"`
}

type CreateConsolidatedShipmentGPResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UpdateConsolidatedShipmentGPRequest struct {
	Interid         string `json:"interid"`
	PRPCSNo         string `json:"prp_cs_no"`
	PrpDriverName   string `json:"prp_driver_name" valid:"required"`
	PrVehicleNumber string `json:"pr_vehicle_number" valid:"required"`
	PhoneName       string `json:"phonname" valid:"required"`

	PurchaseOrders []POLST `json:"polist" valid:"required"`
}

type UpdateConsolidatedShipmentGPResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
