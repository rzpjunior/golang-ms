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
	orm.RegisterModel(new(SummaryRouting))
}

// Routing: struct to hold model data for database
type SummaryRouting struct {
	ID              int64   `orm:"column(id);auto" json:"-"`
	TotalSalesOrder int64   `orm:"column(total_sales_order)" json:"total_sales_order"`
	TotalWeight     float64 `orm:"column(total_weight)" json:"total_weight"`
	TotalKoli       float64 `orm:"column(total_koli)" json:"total_koli"`
	TotalBranch     int64   `orm:"column(total_branch)" json:"total_branch"`
	TotalFragile    float64 `orm:"column(total_fragile)" json:"total_fragile"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SummaryRouting) MarshalJSON() ([]byte, error) {
	type Alias SummaryRouting

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating routing struct into routing table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to routing.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SummaryRouting) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting routing data
// this also will truncated all data from all table
// that have relation with this routing .
func (m *SummaryRouting) Delete() (err error) {
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
func (m *SummaryRouting) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
