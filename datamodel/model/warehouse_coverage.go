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
	orm.RegisterModel(new(WarehouseCoverage))
}

// WarehouseCoverage model for warehouse coverage table.
type WarehouseCoverage struct {
	ID            int64 `orm:"column(id);auto" json:"-"`
	MainWarehouse int8  `orm:"column(main_warehouse);null" json:"main_warehouse"`

	Warehouse         *Warehouse   `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	SubDistrict       *SubDistrict `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`
	ParentWarehouseID *Warehouse   `orm:"column(parent_warehouse_id);null;rel(fk)" json:"parent_warehouse,omitempty"`

	StatusConvert string `orm:"-" json:"status_convert"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *WarehouseCoverage) MarshalJSON() ([]byte, error) {
	type Alias WarehouseCoverage

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *WarehouseCoverage) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *WarehouseCoverage) Delete() (err error) {
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
func (m *WarehouseCoverage) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
