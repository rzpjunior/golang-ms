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
	orm.RegisterModel(new(PackingOrderItem))
}

// PackingOrder model for packing order table.
type PackingOrderItem struct {
	ID          int64   `orm:"column(id);auto" json:"-"`
	TotalOrder  float64 `orm:"column(total_order)" json:"total_order"`
	TotalWeight float64 `orm:"column(total_weight)" json:"total_weight"`
	TotalPack   float64 `orm:"column(total_pack)" json:"total_pack"`

	PackingOrder *PackingOrder `orm:"column(packing_order_id);null;rel(fk)" json:"packing_order,omitempty"`
	Product      *Product      `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
	Helper       []*Staff      `orm:"-" json:"helper,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PackingOrderItem) MarshalJSON() ([]byte, error) {
	type Alias PackingOrderItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PackingOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PackingOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
