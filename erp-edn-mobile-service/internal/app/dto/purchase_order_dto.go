package dto

import "time"

type PurchaseOrderResponse struct {
	ID                     string    `json:"id,omitempty"`
	Code                   string    `json:"code,omitempty"`
	VendorID               string    `json:"vendor_id,omitempty"`
	SiteID                 int64     `json:"site_id,omitempty"`
	TermPaymentPurID       int64     `json:"term_payment_pur_id,omitempty"`
	VendorClassificationID int64     `json:"vendor_classification_id,omitempty"`
	PurchasePlanID         int64     `json:"purchase_plan_id,omitempty"`
	ConsolidatedShipmentID int64     `json:"consolidated_shipment_id,omitempty"`
	StatusGP               int32     `json:"status_gp,omitempty"`
	StatusGPStr            string    `json:"status_gp_str,omitempty"`
	Status                 int32     `json:"status,omitempty"`
	StatusStr              string    `json:"status_str,omitempty"`
	RecognitionDate        time.Time `json:"recognition_date,omitempty"`
	EtaDate                time.Time `json:"eta_date,omitempty"`
	SiteAddress            string    `json:"site_address,omitempty"`
	EtaTime                string    `json:"eta_time,omitempty"`
	TaxPct                 float64   `json:"tax_pct,omitempty"`
	DeliveryFee            float64   `json:"delivery_fee,omitempty"`
	TotalPrice             float64   `json:"total_price,omitempty"`
	TaxAmount              float64   `json:"tax_amount,omitempty"`
	TotalCharge            float64   `json:"total_charge,omitempty"`
	TotalInvoice           float64   `json:"total_invoice,omitempty"`
	TotalWeight            float64   `json:"total_weight,omitempty"`
	Note                   string    `json:"note,omitempty"`
	DeltaPrint             int32     `json:"delta_print,omitempty"`
	Latitude               float64   `json:"latitude,omitempty"`
	Longitude              float64   `json:"longitude,omitempty"`
	UpdatedAt              time.Time `json:"updated_at,omitempty"`
	UpdatedBy              int64     `json:"updated_by,omitempty"`
	CreatedFrom            int32     `json:"created_from,omitempty"`
	HasFinishedGr          int8      `json:"has_finished_gr,omitempty"`
	CreatedAt              time.Time `json:"created_at,omitempty"`
	CreatedBy              int64     `json:"created_by,omitempty"`
	CommittedAt            time.Time `json:"committed_at,omitempty"`
	CommittedBy            int64     `json:"committed_by,omitempty"`
	AssignedTo             int64     `json:"assigned_to,omitempty"`
	AssignedBy             int64     `json:"assigned_by,omitempty"`
	AssignedAt             time.Time `json:"assigned_at,omitempty"`
	Locked                 int32     `json:"locked,omitempty"`
	LockedBy               int64     `json:"locked_by,omitempty"`

	Vendor             *VendorResponse                  `json:"vendor,omitempty"`
	Site               *SiteResponse                    `json:"site,omitempty"`
	PaymentTerm        *PaymentTermResponse             `json:"payment_term,omitempty"`
	PurchaseOrderItems []*PurchaseOrderItemResponse     `json:"purchase_order_items,omitempty"`
	Receiving          []*ReceivingListinDetailResponse `json:"goods_receipt,omitempty"`
	PurchaseInvoice    *PurchaseInvoiceDetailResponse   `json:"purchase_invoice,omitempty"`
}

type CreatePurchaseOrderRequest struct {
	VendorID      int64   `json:"vendor_id" valid:"required"`
	SiteID        int64   `json:"site_id" valid:"required"`
	OrderDate     string  `json:"order_date" valid:"required"`
	StrEtaDate    string  `json:"eta_date" valid:"required"`
	EtaTime       string  `json:"eta_time" valid:"required"`
	DeliveryFee   float64 `json:"delivery_fee"`
	Note          string  `json:"note"`
	TaxPct        float64 `json:"tax_pct"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	DueDate       string  `json:"due_date" valid:"required"`
	PaymentTermID string  `json:"pymtrmid"`

	PurchaseOrderItems []*CreatePurchaseOrderItemRequest `json:"purchase_order_items" valid:"required"`
}

type UpdatePurchaseOrderRequest struct {
	Id          int64   `json:"-"`
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

type PurchaseOrderListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type PurchaseOrderDetailRequest struct {
	Id int64 `json:"id"`
}

type CancelPurchaseOrderRequest struct {
	Id   int64  `json:"-"`
	Note string `json:"note"`
}

type CreatePurchaseOrderGPRequest struct {
	Potype                  int32                              `json:"potype"`
	Ponumber                string                             `json:"ponumber"`
	Docdate                 string                             `json:"docdate"`
	Buyerid                 string                             `json:"buyerid"`
	Vendorid                string                             `json:"vendorid"`
	Curncyid                string                             `json:"curncyid"`
	Deprtmnt                string                             `json:"deprtmnt"`
	Locncode                string                             `json:"locncode"`
	Taxschid                string                             `json:"taxschid"`
	Subtotal                float64                            `json:"subtotal"`
	Trdisamt                float64                            `json:"trdisamt"`
	Frtamnt                 float64                            `json:"frtamnt"`
	Miscamnt                float64                            `json:"miscamnt"`
	Taxamnt                 float64                            `json:"taxamnt"`
	PrpPurchasePlanNo       string                             `json:"prp_purchaseplan_no"`
	PrpPaymentMethod        string                             `json:"prp_payment_method"`
	PrpRegion               string                             `json:"prp_region"`
	PrpEstimatedArrivalDate string                             `json:"prp_estimatedarrival_date"`
	NoteText                string                             `json:"notetext"`
	CSReference             *CSReference                       `json:"cs_reference"`
	Detail                  []CreatePurchaseOrderItemGPRequest `json:"detail"`
}

type CreatePurchaseOrderItemGPRequest struct {
	Ord      int32   `json:"ord"`
	Itemnmbr string  `json:"itemnmbr"`
	Uofm     string  `json:"uofm"`
	Qtyorder float64 `json:"qtyorder"`
	Qtycance float64 `json:"qtycance"`
	Unitcost float64 `json:"unitcost"`
	Notetext string  `json:"notetext"`
}

type CSReference struct {
	PrpCSNo       string `json:"prp_cs_no"`
	PrVehicleNo   string `json:"pr_vehicle_no"`
	PrpDriverName string `json:"prp_driver_name"`
	PhoneName     string `json:"phonname"`
}

type GetPurchaseOrderGPListRequest struct {
	Limit                int32     `query:"limit"`
	Offset               int32     `query:"offset"`
	OrderBy              string    `query:"orderby"`
	Ponumber             string    `query:"ponumber"`
	PonumberLike         string    `query:"ponumberlike"`
	Vendorid             string    `query:"vendorid"`
	Vendname             string    `query:"vendname"`
	ReqDateFrom          time.Time `query:"reqDateFrom"`
	ReqDateTo            time.Time `query:"reqDateTo"`
	Locncode             string    `query:"locncode"`
	Postatus             int       `query:"postatus"`
	Itemnmbr             string    `query:"itemnmbr"`
	EstimatedArrivalDate string    `query:"EstimatedArrivalDate"`
	EstimatedArrivalTime string    `query:"EstimatedArrivalTime"`
}

type CommitPurchaseOrderGPRequest struct {
	Docnumber string `json:"docnumber" valid:"required"`
	Docdate   string `json:"docdate" valid:"required"`
}
