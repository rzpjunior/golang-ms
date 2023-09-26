package dto

import "time"

type PurchasePlanResponse struct {
	ID                   int64     `json:"id"`
	Code                 string    `json:"code"`
	VendorOrganizationID int64     `json:"vendor_organization_id"`
	SiteID               int64     `json:"site_id"`
	RecognitionDate      time.Time `json:"recognition_date"`
	EtaDate              time.Time `json:"eta_date"`
	EtaTime              string    `json:"eta_time"`
	TotalPrice           float64   `json:"total_price"`
	TotalWeight          float64   `json:"total_weight"`
	TotalPurchasePlanQty float64   `json:"total_purchase_plan_qty"`
	TotalPurchaseQty     float64   `json:"total_purchase_qty"`
	Note                 string    `json:"note"`
	Status               int32     `json:"status"`
	AssignedTo           int64     `json:"assigned_to"`
	AssignedBy           int64     `json:"assigned_by"`
	AssignedAt           time.Time `json:"assigned_at"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            int64     `json:"created_by"`
}

type AssignPurchasePlanGPRequest struct {
	Interid             string `json:"interid"`
	PrpPurchaseplanNo   string `json:"prp_purchaseplan_no"`
	PrpPurchaseplanUser string `json:"prp_purchaseplan_user"`
}

type AssignPurchasePlanGPResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
