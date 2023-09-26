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
	orm.RegisterModel(new(DeliveryKoli))
}

// DeliveryKoli model for delivery koli table.
type DeliveryKoli struct {
	ID         int64       `orm:"column(id);auto" json:"-"`
	SalesOrder *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
	Koli       *Koli       `orm:"column(koli_id);null;rel(fk)" json:"koli"`
	Quantity   float64     `orm:"column(quantity)" json:"quantity"`
	Note       string      `orm:"column(note);size(250);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryKoli) MarshalJSON() ([]byte, error) {
	type Alias DeliveryKoli

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating DeliveryKoli struct into DeliveryKoli table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to DeliveryKoli.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DeliveryKoli) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting DeliveryKoli data
// this also will truncated all data from all table
// that have relation with this DeliveryKoli.
func (m *DeliveryKoli) Delete() (err error) {
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
func (m *DeliveryKoli) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
