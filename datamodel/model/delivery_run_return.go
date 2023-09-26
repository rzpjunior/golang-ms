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
	orm.RegisterModel(new(DeliveryRunReturn))
}

// DeliveryRunReturn model for delivery_run_return table.
// could also be called Delivery Run Sheet Item Return
type DeliveryRunReturn struct {
	ID          int64     `orm:"column(id);auto" json:"-"`
	Code        string    `orm:"column(code)" json:"code"`
	TotalPrice  float64   `orm:"column(total_price)" json:"total_price"`
	TotalCharge float64   `orm:"column(total_charge)" json:"total_charge"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`

	DeliveryRunSheetItem  *DeliveryRunSheetItem    `orm:"column(delivery_run_sheet_item_id);rel(fk);null" json:"delivery_run_sheet_item,omitempty"`
	DeliveryRunReturnItem []*DeliveryRunReturnItem `orm:"reverse(many)" json:"delivery_run_return_item,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryRunReturn) MarshalJSON() ([]byte, error) {
	type Alias DeliveryRunReturn

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating delivery_run_return struct into delivery_run_return table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to delivery_run_return.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DeliveryRunReturn) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting delivery_run_return data
// this also will truncated all data from all table
// that have relation with this delivery_run_return.
func (m *DeliveryRunReturn) Delete() (err error) {
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
func (m *DeliveryRunReturn) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
