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
	orm.RegisterModel(new(CourierLog))
}

// CourierLog model for delivery_run_sheet_item table.
type CourierLog struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	Latitude  *float64  `orm:"column(latitude)" json:"latitude"`
	Longitude *float64  `orm:"column(longitude)" json:"longitude"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`

	Courier    *Courier    `orm:"column(courier_id);rel(fk);null" json:"courier,omitempty"`
	SalesOrder *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *CourierLog) MarshalJSON() ([]byte, error) {
	type Alias CourierLog

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating courier_log struct into courier_log table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to courier_log.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *CourierLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting courier_log data
// this also will truncated all data from all table
// that have relation with this courier_log.
func (m *CourierLog) Delete() (err error) {
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
func (m *CourierLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
