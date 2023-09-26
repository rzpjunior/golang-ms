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
	orm.RegisterModel(new(SupplierCommodity))
}

// SupplierCommodity : struct to hold supplier type model data for database
type SupplierCommodity struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	Code      string    `orm:"column(code)" json:"code"`
	Name      string    `orm:"column(name)" json:"name"`
	Note      string    `orm:"column(note)" json:"note"`
	Status    int8      `orm:"column(status)" json:"status"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	UpdatedBy *Staff    `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SupplierCommodity) MarshalJSON() ([]byte, error) {
	type Alias SupplierCommodity

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SupplierCommodity) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SupplierCommodity) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
