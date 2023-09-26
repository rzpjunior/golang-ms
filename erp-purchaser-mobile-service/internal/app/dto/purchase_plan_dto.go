package dto

import "time"

type PurchasePlanResponse struct {
	ID                   string                         `json:"id"`
	Code                 string                         `json:"code"`
	VendorOrganization   *VendorOrganizationResponse    `json:"vendor_organization"`
	Site                 *SiteResponse                  `json:"site"`
	RecognitionDate      time.Time                      `json:"recognition_date"`
	EtaDate              time.Time                      `json:"eta_date"`
	EtaTime              string                         `json:"eta_time"`
	TotalPrice           float64                        `json:"total_price"`
	TotalWeight          float64                        `json:"total_weight"`
	TotalPurchasePlanQty float64                        `json:"total_purchase_plan_qty"`
	TotalPurchaseQty     float64                        `json:"total_purchase_qty"`
	Note                 string                         `json:"note"`
	Status               int32                          `json:"status"`
	AssignedTo           *UserResponse                  `json:"assigned_to"`
	AssignedBy           *UserResponse                  `json:"assigned_by"`
	AssignedAt           time.Time                      `json:"assigned_at"`
	CreatedAt            time.Time                      `json:"created_at"`
	CreatedBy            *UserResponse                  `json:"created_by"`
	PurchasePlanItems    []*PurchasePlanItemResponse    `json:"purchase_plan_items,omitempty"`
	TotalSku             int                            `json:"total_sku"`
	TonnagePurchasePlan  []*TonnagePurchasePlanResponse `json:"tonnage"`
}

type TonnagePurchasePlanResponse struct {
	UomName     string  `orm:"-" json:"uom_name"`
	TotalWeight float64 `orm:"-" json:"total_weight"`
}

type SummaryPurchasePlanResponse struct {
	TotalActive   int64 `json:"total_active"`
	TotalAssigned int64 `json:"total_assigned"`
}

type PurchasePlanRequestAssign struct {
	FieldPurchaserID string `json:"field_purchaser_id" valid:"required"`
	Session          UserResponse
}

type PurchasePlanListRequest struct {
	Limit                      int32     `json:"limit"`
	Offset                     int32     `json:"offset"`
	Status                     int32     `json:"status"`
	Search                     string    `json:"search"`
	OrderBy                    string    `json:"order_by"`
	SiteID                     string    `json:"site_id"`
	FieldPurchaser             string    `json:"field_purchaser"`
	RecognitionDateFrom        time.Time `json:"-"`
	RecognitionDateTo          time.Time `json:"-"`
	RecognitionDateFromString  string    `json:"recognition_date_from"`
	RecognitionDateToString    string    `json:"recognition_date_to"`
	PurchasePlanDateFrom       time.Time `json:"-"`
	PurchasePlanDateTo         time.Time `json:"-"`
	PurchasePlanDateFromString string    `json:"purchaseplan_date_from"`
	PurchasePlanDateToString   string    `json:"purchaseplan_date_to"`
}
