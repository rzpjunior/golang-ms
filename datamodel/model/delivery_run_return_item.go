// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(DeliveryRunReturnItem))
}

// DeliveryRunReturnItem model for delivery_run_return_item table.
// could also be called Delivery Run Sheet Item Return Item
type DeliveryRunReturnItem struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	ReceiveQty        float64 `orm:"column(receive_qty)" json:"receive_qty"`
	ReturnReason      int8    `orm:"column(return_reason)" json:"return_reason"`
	ReturnReasonValue string  `orm:"-" json:"return_reason_value"`
	ReturnEvidence    string  `orm:"column(return_evidence)" json:"return_evidence"`
	Subtotal          float64 `orm:"column(subtotal)" json:"subtotal"`

	DeliveryRunReturn *DeliveryRunReturn `orm:"column(delivery_run_return_id);rel(fk);null" json:"delivery_run_return,omitempty"`
	DeliveryOrderItem *DeliveryOrderItem `orm:"column(delivery_order_item_id);rel(fk);null" json:"delivery_order_item,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryRunReturnItem) MarshalJSON() ([]byte, error) {
	type Alias DeliveryRunReturnItem

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating delivery_run_return_item struct into delivery_run_return_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to delivery_run_return_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DeliveryRunReturnItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting delivery_run_return_item data
// this also will truncated all data from all table
// that have relation with this delivery_run_return_item.
func (m *DeliveryRunReturnItem) Delete() (err error) {
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
func (m *DeliveryRunReturnItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
