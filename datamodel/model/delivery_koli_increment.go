// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(DeliveryKoliIncrement))
}

// DeliveryKoliIncrement model for Delivery Koli Increment table.
type DeliveryKoliIncrement struct {
	ID         int64   `orm:"column(id);auto" json:"-"`
	Increment  float64 `orm:"column(increment)" json:"increment"`
	IsRead     int8    `orm:"column(is_read)" json:"is_read"`
	PrintLabel int8    `orm:"column(print_label)" json:"print_label"`
	TotalKoli  float64 `orm:"-" json:"total_koli"`
	HelperCode string  `orm:"-" json:"helper_code"`

	SalesOrder *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
	Helper     *Staff      `orm:"-" json:"-"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryKoliIncrement) MarshalJSON() ([]byte, error) {
	type Alias DeliveryKoliIncrement

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating delivery koli increment struct into delivery koli increment table.
func (m *DeliveryKoliIncrement) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting delivery koli increment data
func (m *DeliveryKoliIncrement) Delete() (err error) {
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
func (m *DeliveryKoliIncrement) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
