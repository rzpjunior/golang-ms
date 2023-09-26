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
	orm.RegisterModel(new(SupplierGroup))
}

// SupplierGroup model for supplier_commodity_badge_type table.
type SupplierGroup struct {
	ID                int64              `orm:"column(id);auto" json:"-"`
	SupplierCommodity *SupplierCommodity `orm:"column(supplier_commodity_id);null;rel(fk)" json:"supplier_commodity,omitempty"`
	SupplierBadge     *SupplierBadge     `orm:"column(supplier_badge_id);null;rel(fk)" json:"supplier_badge,omitempty"`
	SupplierType      *SupplierType      `orm:"column(supplier_type_id);null;rel(fk)" json:"supplier_type,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SupplierGroup) MarshalJSON() ([]byte, error) {
	type Alias SupplierGroup

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SupplierGroup) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

func (m *SupplierGroup) Delete() (err error) {
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
func (m *SupplierGroup) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
