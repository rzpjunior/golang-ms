package dto

import "time"

type ConsolidatedShipmentResponse struct {
	ID                             int64                                    `json:"id"`
	Code                           string                                   `json:"code"`
	DriverName                     string                                   `json:"driver_name"`
	VehicleNumber                  string                                   `json:"vehicle_number"`
	DriverPhoneNumber              string                                   `json:"driver_phone_number"`
	DeltaPrint                     int8                                     `json:"delta_print"`
	Status                         int8                                     `json:"status"`
	CreatedAt                      time.Time                                `json:"created_at"`
	CreatedBy                      *UserResponse                            `json:"created_by"`
	ConsolidatedShipmentSignatures []*ConsolidatedShipmentSignatureResponse `json:"consolidated_shipment_signatures,omitempty"`
	SiteName                       string                                   `json:"site_name"`
	PurchaseOrders                 []*PurchaseOrderResponse                 `json:"purchase_orders,omitempty"`
	SkuSummaries                   []*SkuSummaryResponse                    `json:"sku_summaries"`
}

type SkuSummaryResponse struct {
	ItemID         string                     `json:"-"`
	ItemName       string                     `json:"item_name"`
	UomName        string                     `json:"uom_name"`
	TotalQty       float64                    `json:"total_qty"`
	PurchaseOrders []*PurchaseOrderSKUSummary `json:"purchase_orders"`
}

type PurchaseOrderSKUSummary struct {
	PurchaseOrderCode string  `orm:"-" json:"purchase_order_code"`
	Qty               float64 `orm:"-" json:"qty"`
}

type PurchaseOrderConsolidateShipmentResponse struct {
	PurchaseOrderCode string  `json:"purchase_order_code"`
	Qty               float64 `json:"qty"`
}

type ConsolidatedShipmentRequestList struct {
	Limit         int32     `json:"limit"`
	Offset        int32     `json:"offset"`
	Status        int32     `json:"status"`
	Search        string    `json:"search"`
	OrderBy       string    `json:"order_by"`
	SiteID        string    `json:"site_id"`
	CreatedAtFrom time.Time `json:"created_at_from"`
	CreatedAtTo   time.Time `json:"created_at_to"`
	CreatedBy     int64     `json:"-"`
}

type ConsolidatedShipmentRequestCreate struct {
	CSNo              string                                     `json:"cs_no"`
	DriverName        string                                     `json:"driver_name" valid:"required|alpha_num_space|lte:100|gte:4"`
	VehicleNumber     string                                     `json:"vehicle_number" valid:"required|alpha_num|range:3,9"`
	DriverPhoneNumber string                                     `json:"driver_phone_number" valid:"required|numeric|range:8,15"`
	PurchaseOrders    []*PurchaseOrderConsolidateShipmentRequest `json:"purchase_orders" valid:"required"`
}

type PurchaseOrderConsolidateShipmentRequest struct {
	PurchaseOrderID string `json:"purchase_order_id" valid:"required"`
}

type ConsolidatedShipmentRequestUpdate struct {
	CSNo              string                                     `json:"cs_no" valid:"required"`
	DriverName        string                                     `json:"driver_name" valid:"required|alpha_num_space|lte:100|gte:4"`
	VehicleNumber     string                                     `json:"vehicle_number" valid:"required|alpha_num|range:5,9"`
	DriverPhoneNumber string                                     `json:"driver_phone_number" valid:"required|numeric|range:8,15"`
	PurchaseOrders    []*PurchaseOrderConsolidateShipmentRequest `json:"purchase_orders" valid:"required"`
}

type ConsolidatedShipmentRequestSignature struct {
}
