// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// DeliveryRunSheetItem model for delivery_run_sheet_item table.
type DeliveryRunSheetItem struct {
	ID                          int64     `orm:"column(id)" json:"-"`
	StepType                    int8      `orm:"column(step_type)" json:"step_type"`
	Latitude                    *float64  `orm:"column(latitude)" json:"latitude"`
	Longitude                   *float64  `orm:"column(longitude)" json:"longitude"`
	Status                      int8      `orm:"column(status);null" json:"status"`
	Note                        string    `orm:"column(note)" json:"note"`
	RecipientName               string    `orm:"column(recipient_name)" json:"recipient_name"`
	MoneyReceived               float64   `orm:"column(money_received)" json:"money_received"`
	DeliveryEvidenceImageURL    string    `orm:"column(delivery_evidence_image_url)" json:"delivery_evidence_image_url"`
	TransactionEvidenceImageURL string    `orm:"column(transaction_evidence_image_url)" json:"transaction_evidence_image_url"`
	ArrivalTime                 time.Time `orm:"column(arrival_time)" json:"arrival_time"`
	UnpunctualReason            int8      `orm:"column(unpunctual_reason)" json:"unpunctual_reason"`
	UnpunctualDetail            int8      `orm:"column(unpunctual_detail)" json:"unpunctual_detail"`
	FarDeliveryReason           string    `orm:"column(far_delivery_reason)" json:"far_delivery_reason"`
	CreatedAt                   time.Time `orm:"column(created_at)" json:"created_at"`
	StartedAt                   time.Time `orm:"column(started_at)" json:"started_at"`
	FinishedAt                  time.Time `orm:"column(finished_at)" json:"finished_at"`

	DeliveryRunSheetID int64  `orm:"column(delivery_run_sheet_id)" json:"delivery_run_sheet_id"`
	CourierID          string `orm:"column(courier_id)" json:"courier_id"`
	SalesOrderID       string `orm:"column(sales_order_id)" json:"sales_order_id"`
}

func init() {
	orm.RegisterModel(new(DeliveryRunSheetItem))
}

func (m *DeliveryRunSheetItem) TableName() string {
	return "delivery_run_sheet_item"
}
