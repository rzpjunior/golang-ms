// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type DeliveryRunSheetItemGetRequest struct {
	Offset               int
	Limit                int
	OrderBy              string
	GroupBy              string
	StepType             []int
	Status               []int
	DeliveryRunSheetIDs  []int64
	CourierIDs           []string
	ArrSalesOrderIDs     []string
	SearchSalesOrderCode string
}

type GroupedDeliveryRunSheetItemGetRequest struct {
	Offset                      int
	Limit                       int
	OrderBy                     string
	Status                      []int
	ArrSalesOrderIDs            []string
	ArrCourierVendorsCourierIDs []string
	CourierID                   string
	SearchSalesOrderID          string
}

type DeliveryRunSheetItemResponse struct {
	ID                          int64     `json:"id"`
	StepType                    int8      `json:"step_type"`
	Latitude                    *float64  `json:"latitude"`
	Longitude                   *float64  `json:"longitude"`
	Status                      int8      `json:"status"`
	Note                        string    `json:"note"`
	RecipientName               string    `json:"recipient_name"`
	MoneyReceived               float64   `json:"money_received"`
	DeliveryEvidenceImageURL    string    `json:"delivery_evidence_image_url"`
	TransactionEvidenceImageURL string    `json:"transaction_evidence_image_url"`
	ArrivalTime                 time.Time `json:"arrival_time"`
	UnpunctualReason            int8      `json:"unpunctual_reason"`
	UnpunctualDetail            int8      `json:"unpunctual_detail"`
	FarDeliveryReason           string    `json:"far_delivery_reason"`
	CreatedAt                   time.Time `json:"created_at"`
	StartedAt                   time.Time `json:"started_at"`
	FinishedAt                  time.Time `json:"finished_at"`

	DeliveryRunSheetID int64  `json:"delivery_run_sheet_id"`
	CourierID          string `json:"courier_id"`
	SalesOrderID       string `json:"sales_order_id"`
}
