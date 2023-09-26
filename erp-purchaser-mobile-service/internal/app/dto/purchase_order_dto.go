package dto

import "time"

type PurchaseOrderResponse struct {
	ID                     string                            `json:"id"`
	Code                   string                            `json:"code"`
	Vendor                 *VendorResponse                   `json:"vendor"`
	Site                   *SiteResponse                     `json:"site"`
	TermPaymentPur         *PurchaseTermResponse             `json:"term_payment_pur"`
	VendorClassification   *VendorClassificationResponse     `json:"vendor_classification"`
	PurchasePlan           *PurchasePlanResponse             `json:"purchase_plan"`
	ConsolidatedShipment   *ConsolidatedShipmentResponse     `json:"consolidated_shipment,omitempty"`
	Status                 int32                             `json:"status"`
	DocDate                time.Time                         `json:"doc_date"`
	RecognitionDate        time.Time                         `json:"recognition_date"`
	EtaDate                time.Time                         `json:"eta_date"`
	SiteAddress            string                            `json:"site_address"`
	EtaTime                string                            `json:"eta_time"`
	TaxPct                 float64                           `json:"tax_pct"`
	DeliveryFee            float64                           `json:"delivery_fee"`
	TotalPrice             float64                           `json:"total_price"`
	TaxAmount              float64                           `json:"tax_amount"`
	TotalCharge            float64                           `json:"total_charge"`
	TotalInvoice           float64                           `json:"total_invoice"`
	TotalWeight            float64                           `json:"total_weight"`
	Note                   string                            `json:"note"`
	DeltaPrint             int32                             `json:"delta_print"`
	Latitude               float64                           `json:"latitude"`
	Longitude              float64                           `json:"longitude"`
	CreatedFrom            int32                             `json:"created_from"`
	HasFinishedGr          int32                             `json:"has_finished_gr"`
	CreatedAt              time.Time                         `json:"created_at"`
	CreatedBy              *UserResponse                     `json:"created_by"`
	CommittedAt            time.Time                         `json:"committed_at"`
	CommittedBy            *UserResponse                     `json:"committed_by"`
	AssignedAt             time.Time                         `json:"assigned_at"`
	AssignedTo             *UserResponse                     `json:"assigned_to"`
	AssignedBy             *UserResponse                     `json:"assigned_by"`
	UpdatedAt              time.Time                         `json:"updated_at"`
	UpdatedBy              *UserResponse                     `json:"updated_by"`
	Locked                 int32                             `json:"locked"`
	LockedBy               *UserResponse                     `json:"locked_by"`
	PurchaseOrderItems     []*PurchaseOrderItemResponse      `json:"purchase_order_items,omitempty"`
	TotalSku               int                               `json:"total_sku"`
	TonnagePurchaseOrder   []*TonnagePurchaseOrderResponse   `json:"tonnage"`
	PurchaseOrderSignature []*PurchaseOrderSignatureResponse `json:"signature,omitempty"`
	PurchaseOrderImage     []*PurchaseOrderImageResponse     `json:"images,omitempty"`
}

type TonnagePurchaseOrderResponse struct {
	UomName     string  `orm:"-" json:"uom_name"`
	TotalWeight float64 `orm:"-" json:"total_weight"`
}

type PurchaseOrderSignatureResponse struct {
	ID           int64     `json:"id"`
	JobFunction  string    `json:"job_function"`
	Name         string    `json:"name"`
	SignatureURL string    `json:"signature_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type PurchaseOrderImageResponse struct {
	ID        int64     `json:"id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type PurchaseOrderListRequest struct {
	Limit               int32     `json:"limit"`
	Offset              int32     `json:"offset"`
	Status              int32     `json:"status"`
	Search              string    `json:"search"`
	OrderBy             string    `json:"order_by"`
	RecognitionDateFrom time.Time `json:"recognition_date_from"`
	RecognitionDateTo   time.Time `json:"recognition_date_to"`
	Code                []string  `json:"code"`
	IsNotConsolidated   bool
	PurchasePlanID      string `json:"purchase_plan_id"`
	EmployeeCode        string `json:"employee_code"`
	Site                string `json:"site"`
	PrpCsNo             string `json:"prp_cs_no"`
}

type PurchaseOrderRequestCreate struct {
	PurchasePlanID     string                            `json:"purchase_plan_id" valid:"required"`
	VendorID           string                            `json:"vendor_id" valid:"required"`
	PurchaseTermID     string                            `json:"term_payment_pur_id"`
	RegionID           string                            `json:"region_id" valid:"required"`
	SiteID             string                            `json:"site_id" valid:"required"`
	SiteAddress        string                            `json:"site_address" valid:"required"`
	RecognitionDate    string                            `json:"order_date" valid:"required"`
	EtaDate            string                            `json:"eta_date" valid:"required"`
	EtaTime            string                            `json:"eta_time" valid:"required"`
	DeliveryFee        float64                           `json:"delivery_fee"`
	Note               string                            `json:"note"`
	TaxPct             float64                           `json:"tax_pct"`
	CreatedFrom        int8                              `json:"created_from"`
	Latitude           float64                           `json:"latitude"`
	Longitude          float64                           `json:"longitude"`
	PurchaseOrderItems []*PurchaseOrderItemRequestCreate `json:"purchase_order_items" valid:"required"`
	Images             []string                          `json:"images" valid:"required"`
	DueDate            string                            `json:"duedate"`
	PaymentTermID      string                            `json:"term_payment_id"`
	PRStatus           int32                             `json:"pr_status"`
}

type PurchaseOrderRequestUpdate struct {
	PurchaseOrderItems []*PurchaseOrderItemRequestCreate `json:"purchase_order_items" valid:"required"`
	Images             []string                          `json:"images" valid:"required"`
}

type PurchaseOrderRequestAssign struct {
	FieldPurchaserID int64 `json:"field_purchaser_id" valid:"required"`
}

type PurchaseOrderRequestCancel struct {
	FieldPurchaserID string `json:"field_purchaser_id"`
}
type PurchaseOrderRequestSignature struct {
	PurchaseOrderID string `json:"purchase_order_id" valid:"required"`
	JobFunction     string `json:"job_function" valid:"required|alpha_num_space|lte:100"`
	Name            string `json:"name" valid:"required|alpha_num_space|lte:100"`
	SignatureURL    string `json:"signature_url" valid:"required"`
}
