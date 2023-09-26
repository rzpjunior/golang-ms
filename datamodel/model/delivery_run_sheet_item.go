// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(DeliveryRunSheetItem))
}

// DeliveryRunSheetItem model for delivery_run_sheet_item table.
type DeliveryRunSheetItem struct {
	ID                          int64     `orm:"column(id);auto" json:"-"`
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
	UnpunctualReasonValue       string    `orm:"-" json:"unpunctual_reason_value"`
	FarDeliveryReason           string    `orm:"column(far_delivery_reason)" json:"far_delivery_reason"`
	CreatedAt                   time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	StartedAt                   time.Time `orm:"column(started_at);type(timestamp);null" json:"started_at"`
	FinishedAt                  time.Time `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`

	// for delivery run sheet detail
	CustomerLatitude    *float64 `orm:"-" json:"customer_latitude"`
	CustomerLongitude   *float64 `orm:"-" json:"customer_longitude"`
	TotalSalesOrder     int64    `orm:"-" json:"total_sales_order"`
	CompletedSalesOrder int64    `orm:"-" json:"completed_sales_order"`
	Distance            float64  `orm:"-" json:"distance"`
	DistanceUnit        string   `orm:"-" json:"distance_unit"`

	DeliveryRunSheet *DeliveryRunSheet `orm:"column(delivery_run_sheet_id);rel(fk);null" json:"delivery_run_sheet,omitempty"`
	Courier          *Courier          `orm:"column(courier_id);rel(fk);null" json:"courier,omitempty"`
	SalesOrder       *SalesOrder       `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order,omitempty"`

	DeliveryRunReturn   *DeliveryRunReturn     `orm:"-" json:"delivery_run_return"`
	PostponeDeliveryLog []*PostponeDeliveryLog `orm:"-" json:"postpone_delivery_log"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryRunSheetItem) MarshalJSON() ([]byte, error) {
	type Alias DeliveryRunSheetItem

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating delivery_run_sheet_item struct into delivery_run_sheet_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to delivery_run_sheet_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DeliveryRunSheetItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting delivery_run_sheet_item data
// this also will truncated all data from all table
// that have relation with this delivery_run_sheet_item.
func (m *DeliveryRunSheetItem) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *DeliveryRunSheetItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
