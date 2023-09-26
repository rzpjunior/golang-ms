// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

// GET DRS
type ControlTowerGetDRSRequest struct {
	Offset            int
	Limit             int
	OrderBy           string
	SiteID            string
	StartDeliveryDate time.Time
	EndDeliveryDate   time.Time
	CourierVendorID   string
	CourierID         string
	StatusIDs         []int
	Search            string
}
type ControlTowerGetDRSResponse struct {
	CompletedSalesOrder int64                                  `json:"completed_sales_order"`
	TotalSalesOrder     int64                                  `json:"total_sales_order"`
	Courier             *SubControlTowerGetDRSCourier          `json:"courier"`
	DeliveryRunSheet    *SubControlTowerGetDRSDeliveryRunSheet `json:"delivery_run_sheet"`
}
type SubControlTowerGetDRSCourier struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type SubControlTowerGetDRSDeliveryRunSheet struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	DeliveryDate time.Time `json:"delivery_date"`
	Status       int8      `json:"status"`
}

// GET COURIER
type ControlTowerGetCourierRequest struct {
	SiteID          string `json:"site_id"`
	CourierVendorID string `json:"courier_vendor_id"`
	CourierID       string `json:"courier_id"`
}
type ControlTowerGetCourierResponse struct {
	Couriers    []*SubControlTowerGetCourierCourier `json:"courier"`
	OnEmergency int64                               `json:"on_emergency"`
}
type SubControlTowerGetCourierCourier struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PhoneNumber   string    `json:"phone_number"`
	LicensePlate  string    `json:"license_plate"`
	EmergencyMode int8      `json:"emergency_mode"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	LastUpdated   time.Time `json:"last_updated"`

	// vehicle profile
	VehicleProfileType string `json:"vehicle_profile_type"`
	VendorCourierCode  string `json:"vendor_courier_code"`
}

// GET DETAIL DRS
type ControlTowerGetDRSDetailResponse struct {
	Id                int64                               `json:"id"`
	Code              string                              `json:"code"`
	DeliveryDate      time.Time                           `json:"delivery_date"`
	Courier           *SubControlTowerGetDRSDetailCourier `json:"courier"`
	StartedAt         time.Time                           `json:"started_at"`
	FinishedAt        time.Time                           `json:"finished_at"`
	StartingLatitude  *float64                            `json:"starting_latitude"`
	StartingLongitude *float64                            `json:"starting_longitude"`
	FinishedLatitude  *float64                            `json:"finished_latitude"`
	FinishedLongitude *float64                            `json:"finished_longitude"`
	Status            int8                                `json:"status"`
}
type SubControlTowerGetDRSDetailCourier struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	CourierPhoneNumber   string `json:"courier_phone_number"`
	CourierVendorName    string `json:"courier_vendor_name"`
	CourierVehicleName   string `json:"vehicle_profile_name"`
	LicensePlate         string `json:"license_plate"`
	CourierWarehouseName string `json:"warehouse_name"`
}

// GET DETAIL COURIER
type ControlTowerGetCourierDetailResponse struct {
	TotalSalesOrder      int64                                   `json:"total_sales_order"`
	TotalSelfPickup      int64                                   `json:"total_self_pickup"`
	TotalDeliveryReturn  int64                                   `json:"total_delivery_return"`
	Courier              *SubControlTowerGetCourierDetailCourier `json:"courier"`
	DeliveryRunSheetItem []*SubControlTowerGetCourierDetailDRSI  `json:"delivery_run_sheet_item"`
}
type SubControlTowerGetCourierDetailCourier struct {
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	EmergencyMode     int8      `json:"emergency_mode"`
	LastEmergencyTime time.Time `json:"last_emergency_time"`
	LastUpdated       time.Time `json:"last_updated"`

	// vehicle profile
	VehicleProfileType string `json:"vehicle_profile_type"`
}
type SubControlTowerGetCourierDetailDRSI struct {
	Id                          int64     `json:"id"`
	RecipientName               string    `json:"recipient_name"`
	MoneyReceived               float64   `json:"money_received"`
	StartTime                   time.Time `json:"start_time"`
	ArrivalTime                 time.Time `json:"arrival_time"`
	FinishTime                  time.Time `json:"finish_time"`
	UnpunctualDetail            int8      `json:"unpunctual_detail"`
	UnpunctualReason            int8      `json:"unpunctual_reason"`
	UnpunctualReasonValue       string    `json:"unpunctual_reason_value"`
	FarDeliveryReason           string    `json:"far_delivery_reason"`
	DeliveryEvidenceImageURL    string    `json:"delivery_evidence_image_url"`
	TransactionEvidenceImageURL string    `json:"transaction_evidence_image_url"`
	Note                        string    `json:"note"`
	Status                      int8      `json:"status"`

	SalesOrder        *SubControlTowerGetCourierDetailSO                    `json:"sales_order"`
	DeliveryRunReturn *SubControlTowerGetCourierDetailDRR                   `json:"delivery_run_return"`
	PostponeLog       []*SubControlTowerGetCourierDetailPostponeDeliveryLog `json:"postpone_delivery_log"`
}
type SubControlTowerGetCourierDetailSO struct {
	ID                    int64   `json:"id"`
	Code                  string  `json:"code"`
	DeliveryDate          string  `json:"delivery_date"`
	DeliveryFee           float64 `json:"delivery_fee"`
	VoucherDiscountAmount float64 `json:"vou_disc_amount"`
	PointRedeemAmount     float64 `json:"point_redeem_amount"`

	// branch coordinate log
	CustomerLatitude  float64 `json:"customer_latitude"`
	CustomerLongitude float64 `json:"customer_longitude"`

	// customer
	CustomerName string `json:"customer_name"`

	// address
	AddressName        string `json:"address_name"`
	AddressPhoneNumber string `json:"address_phone_number"`
	ShippingAddress    string `json:"shipping_address"`
	//// subdistrict
	SubDistrictDetail string `json:"sub_district_detail"`
	PostalCode        string `json:"postal_code"`

	// wrt
	WrtName string `json:"wrt_name"`

	// term payment sls
	PaymentTypeName string `json:"payment_type_name"`

	// sales invoice
	SalesInvoice *SalesInvoice `json:"sales_invoice"`
}

type SalesInvoice struct {
	TotalCharge float64 `json:"total_charge"`
}
type SubControlTowerGetCourierDetailDRR struct {
	TotalPrice            float64                                `json:"total_price,omitempty"`
	TotalCharge           float64                                `json:"total_charge,omitempty"`
	DeliveryRunReturnItem []*SubControlTowerGetCourierDetailDRRI `json:"delivery_run_return_item,omitempty"`
}
type SubControlTowerGetCourierDetailDRRI struct {
	DeliveryQty       float64 `json:"delivery_qty"`
	ReceiveQty        float64 `json:"receive_qty"`
	Subtotal          float64 `json:"subtotal"`
	ReturnReason      int8    `json:"return_reason"`
	ReturnReasonValue string  `json:"return_reason_value"`
	ReturnEvidence    string  `json:"return_evidence"`

	// product
	ProductName string `json:"product_name"`
	//// Unit of Measurement (U of M)
	UOMName string `json:"uom_name"`
}
type SubControlTowerGetCourierDetailPostponeDeliveryLog struct {
	PostponeReason   string    `json:"postpone_reason"`
	PostponeEvidence string    `json:"postpone_evidence"`
	StartedAt        time.Time `json:"started_at"`
	PostponedAt      time.Time `json:"postponed_at"`
}

// CANCEL DRS
type ControlTowerCancelDRSRequest struct {
	DeliveryRunSheetID int64  `json:"-"`
	Note               string `json:"note" valid:"required"`
}
type ControlTowerCancelDRSResponse struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Status int8   `json:"status"`
}

// CANCEL ITEM
type ControlTowerCancelItemRequest struct {
	DeliveryRunSheetItemID int64  `json:"-"`
	Note                   string `json:"note" valid:"required"`
}
type ControlTowerCancelItemResponse struct {
	ID     int64 `json:"id"`
	Status int8  `json:"status"`
}
