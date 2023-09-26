// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

// Login
type LoginRequest struct {
	CourierCode string `json:"code" valid:"required"`
	Password    string `json:"password" valid:"required"`
	Timezone    string `json:"timezone"`
}
type LoginResponse struct {
	Courier *GlobalCourier `json:"courier,omitempty"`
	Token   string         `json:"token"`
}

// Create Courier Log
type CreateCourierLogRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	CourierID string `json:"-"`
}

// Get
type CourierAppGetRequest struct {
	Offset               int
	Limit                int
	OrderBy              string
	CourierId            string
	StartDeliveryDate    time.Time
	EndDeliveryDate      time.Time
	StepType             int
	StatusIDs            []int
	Search               string
	SearchSalesOrderCode string
}
type CourierAppGetResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
	Distance             float64                     `json:"distance"`
}

// Get Detail
type CourierAppDetailResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
	Distance             float64                     `json:"distance"`

	// SalesOrder *GlobalSalesOrder `json:"sales_order"`
	SalesOrder        *SubCourierAppDetailSO   `json:"sales_order"`
	DeliveryRunReturn *GlobalDeliveryRunReturn `json:"delivery_run_return"`
	Koli              int64                    `json:"koli"`
}
type SubCourierAppDetailSO struct {
	Code                  string  `json:"code"`
	DeliveryDate          string  `json:"delivery_date"`
	DeliveryFee           float64 `json:"delivery_fee"`
	VoucherDiscountAmount float64 `json:"vou_disc_amount"`
	PointRedeemAmount     float64 `json:"point_redeem_amount"`
	SalesPayment          string  `json:"sales_payment"`
	// sales invoice
	SalesInvoiceTotalCharge float64 `json:"sales_invoice_total_charge"`

	Customer  *GlobalCustomer         `json:"customer"`
	Address   *GlobalAddress          `json:"address"`
	Wrt       *GlobalWrt              `json:"wrt"`
	OrderType *GlobalOrderType        `json:"order_type"`
	Item      []*GlobalSalesOrderItem `json:"sales_order_item"`
}

// Scan Detail
type CourierAppScanDetailRequest struct {
	Code string `json:"code" valid:"required"`

	CourierId string `json:"-"`
}

// Scan
type CourierAppScanRequest struct {
	Code string `json:"code" valid:"required"`

	CourierSiteId string `json:"-"`
}
type CourierAppScanResponse struct {
	SalesOrder *GlobalSalesOrder `json:"sales_order"`
	Koli       int64             `json:"koli"`
}

// Self Assign
type CourierAppSelfAssignRequest struct {
	SopNumber string  `json:"sop_number" valid:"required"`
	Latitude  float64 `json:"latitude" valid:"required"`
	Longitude float64 `json:"longitude" valid:"required"`

	CourierId     string `json:"-"`
	CourierSiteId string `json:"-"`
}

type CourierAppSelfAssignResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
}

// Start Delivery
type CourierAppStartDeliveryRequest struct {
	Id        int64   `json:"-" valid:"required"`
	Latitude  float64 `json:"latitude" valid:"required"`
	Longitude float64 `json:"longitude" valid:"required"`

	CourierId     string `json:"-"`
	CourierSiteId string `json:"-"`
}
type CourierAppStartDeliveryResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
}

// Success Delivery
type CourierAppSuccessDeliveryRequest struct {
	Id                  int64   `json:"-" valid:"required"` // ID Delivery run sheet item
	Latitude            float64 `json:"latitude" valid:"required"`
	Longitude           float64 `json:"longitude" valid:"required"`
	RecipientName       string  `json:"recipient_name" valid:"required"`
	MoneyReceived       float64 `json:"money_received"`
	DeliveryEvidence    string  `json:"delivery_evidence" valid:"required"`
	TransactionEvidence string  `json:"transaction_evidence"`
	UnpunctualReason    int8    `json:"unpunctual_reason"`
	FarDeliveryReason   string  `json:"far_delivery_reason"`
	Note                string  `json:"note"`

	UnpunctualDetail int8 `json:"-"`

	CourierId     string `json:"-"`
	CourierSiteId string `json:"-"`
}
type CourierAppSuccessDeliveryResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
}

// Postpone Delivery
type CourierAppPostponeDeliveryRequest struct {
	Id                       int64  `json:"-" valid:"required"` // ID Delivery run sheet item
	Note                     string `json:"note" valid:"required"`
	PostponeDeliveryEvidence string `json:"postpone_delivery_evidence" valid:"required"`

	CourierId string `json:"-"`
}
type CourierAppPostponeDeliveryResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
}

// Fail Delivery
type CourierAppFailDeliveryRequest struct {
	Id        int64   `json:"-" valid:"required"` // ID Delivery run sheet item
	Latitude  float64 `json:"latitude" valid:"required"`
	Longitude float64 `json:"longitude" valid:"required"`
	Note      string  `json:"note" valid:"required"`

	CourierId string `json:"-"`
}
type CourierAppFailDeliveryResponse struct {
	DeliveryRunSheetItem *GlobalDeliveryRunSheetItem `json:"delivery_run_sheet_item"`
}

// Status Delivery
type CourierAppStatusDeliveryRequest struct {
	Id        int64   `json:"-" valid:"required"` // ID Delivery run sheet item
	Latitude  float64 `json:"latitude" valid:"required"`
	Longitude float64 `json:"longitude" valid:"required"`
}
type CourierAppStatusDeliveryResponse struct {
	Punctual bool `json:"punctual"` // false = unpuctual , true = punctual
	Earlier  bool `json:"earlier"`  // false = late , true = early
	Nearby   bool `json:"nearby"`   // false = far , true = near
}

// Activate Emergency
type CourierAppActivateEmergencyRequest struct {
	CourierId string `json:"-"`
}
type CourierAppActivateEmergencyResponse struct {
	Courier *GlobalCourier `json:"courier"`
}

// Deactivate Emergency
type CourierAppDeactivateEmergencyRequest struct {
	CourierId string `json:"-"`
}
type CourierAppDeactivateEmergencyResponse struct {
	Courier *GlobalCourier `json:"courier"`
}

// Create Merchant Delivery Log
type CourierAppCreateMerchantDeliveryLogRequest struct {
	Id        int64   `json:"-" valid:"required"` // ID Delivery run sheet item
	Latitude  float64 `json:"latitude" valid:"required"`
	Longitude float64 `json:"longitude" valid:"required"`

	CourierId string `json:"-"`
}
type CourierAppCreateMerchantDeliveryLogResponse struct {
	MerchantDeliveryLog *GlobalMerchantDeliveryLog `json:"merchant_delivery_log"`
}

// Get Glossary
type CourierAppGetGlossaryRequest struct {
	Table     string
	Attribute string
	ValueInt  int
	ValueName string
}

type CourierAppGetGlossaryResponse struct {
	Glossary *GlobalGlossary `json:"glossary"`
}

type DeliveryReturnRequest struct {
	ID                int64                           `json:"-" valid:"required"` // ID Delivery run sheet item
	CourierID         string                          `json:"-"`
	TotalPrice        float64                         `json:"-"`
	GrandTotal        float64                         `json:"-"`
	Code              string                          `json:"-"`
	Items             []*DeliveryRunReturnItemRequest `json:"items" valid:"required"`
	ReturnedSomething bool                            `json:"-"` // false = none , yes = returned something

}

type DeliveryRunReturnItemRequest struct {
	ItemNumber              string  `json:"item_number"`
	DeliveryRunReturnItemID int64   `json:"delivery_run_return_item_id"`
	DeliveryQty             float64 `json:"delivery_qty"`
	ReceiveQty              float64 `json:"receive_qty"`
	ReturnReason            int8    `json:"return_reason"`
	ReturnEvidence          string  `json:"return_evidence"`
	Subtotal                float64 `json:"-"`
}

type DeleteDeliveryReturnRequest struct {
	ID        int64  `json:"-" valid:"required"` // ID Delivery run sheet item
	CourierID string `json:"-"`
}
