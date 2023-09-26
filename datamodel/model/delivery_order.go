// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(DeliveryOrder))
}

// DeliveryOrder: struct to hold model data for database
type DeliveryOrder struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate  time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	ShippingAddress  string    `orm:"column(shipping_address)" json:"shipping_address"`
	ReceiptNote      string    `orm:"column(receipt_note)" json:"receipt_note"`
	TotalWeight      float64   `orm:"column(total_weight)" json:"total_weight"`
	DeltaPrint       int8      `orm:"column(delta_print)" json:"delta_print"`
	Note             string    `orm:"column(note)" json:"note"`
	Status           int8      `orm:"column(status);null" json:"status"`
	CancellationNote string    `orm:"column(cancellation_note);null" json:"cancellation_note,omitempty"`
	HasDelivered     int8      `orm:"column(has_delivered);null" json:"has_delivered"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy        int64     `orm:"column(created_by)" json:"created_by"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	UpdatedBy        int64     `orm:"column(updated_by)" json:"updated_by"`
	ConfirmedAt      time.Time `orm:"column(confirmed_at);type(timestamp);null" json:"confirmed_at"`
	ConfirmedBy      int64     `orm:"column(confirmed_by)" json:"confirmed_by"`
	CancelledAt      time.Time `orm:"column(cancelled_at);type(timestamp);null" json:"cancelled_at"`
	CancelledBy      int64     `orm:"column(cancelled_by)" json:"cancelled_by"`

	SalesOrder *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
	Warehouse  *Warehouse  `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	Wrt        *Wrt        `orm:"column(wrt_id);null;rel(fk)" json:"wrt"`

	DeliveryOrderItems []*DeliveryOrderItem `orm:"reverse(many)" json:"delivery_order_items,omitempty"`
	SalesInvoice       []*SalesInvoice      `orm:"-" json:"sales_invoice,omitempty"`
	DeliveryKoli       []*DeliveryKoli      `orm:"-" json:"delivery_koli"`
	TotalKoli          float64              `orm:"-" json:"total_koli"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DeliveryOrder) MarshalJSON() ([]byte, error) {
	type Alias DeliveryOrder

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusDoc(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DeliveryOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DeliveryOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
