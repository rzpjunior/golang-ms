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
	orm.RegisterModel(new(DeliveryRunSheet))
}

// DeliveryRunSheet model for delivery_run_sheet table.
type DeliveryRunSheet struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	Code              string    `orm:"column(code)" json:"code"`
	DeliveryDate      time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	StartedAt         time.Time `orm:"column(started_at);type(timestamp);null" json:"started_at"`
	FinishedAt        time.Time `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
	StartingLatitude  *float64  `orm:"column(starting_latitude)" json:"starting_latitude"`
	StartingLongitude *float64  `orm:"column(starting_longitude)" json:"starting_longitude"`
	FinishedLatitude  *float64  `orm:"column(finished_latitude)" json:"finished_latitude"`
	FinishedLongitude *float64  `orm:"column(finished_longitude)" json:"finished_longitude"`
	Status            int8      `orm:"column(status);null" json:"status"`

	Courier              *Courier                `orm:"column(courier_id);rel(fk);null" json:"courier,omitempty"`
	DeliveryRunSheetItem []*DeliveryRunSheetItem `orm:"reverse(many)" json:"delivery_run_sheet_item,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryRunSheet) MarshalJSON() ([]byte, error) {
	type Alias DeliveryRunSheet

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating delivery_run_sheet struct into delivery_run_sheet table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to delivery_run_sheet.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DeliveryRunSheet) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting delivery_run_sheet data
// this also will truncated all data from all table
// that have relation with this delivery_run_sheet.
func (m *DeliveryRunSheet) Delete() (err error) {
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
func (m *DeliveryRunSheet) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
