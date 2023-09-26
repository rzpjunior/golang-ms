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
	orm.RegisterModel(new(PackingOrderItemAssign))
}

// PackingOrder model for packing order table.
type PackingOrderItemAssign struct {
	ID             int64   `orm:"column(id);auto" json:"-"`
	SubTotalWeight float64 `orm:"column(subtotal_weight)" json:"subtotal_weight"`
	SubTotalPack   float64 `orm:"column(subtotal_pack)" json:"subtotal_pack"`

	PackingOrderItem *PackingOrderItem `orm:"column(packing_order_item_id);null;rel(fk)" json:"packing_order_item,omitempty"`
	Staff            *Staff            `orm:"column(staff_id);null;rel(fk)" json:"staff,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PackingOrderItemAssign) MarshalJSON() ([]byte, error) {
	type Alias PackingOrderItemAssign

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PackingOrderItemAssign) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PackingOrderItemAssign) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
