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
	orm.RegisterModel(new(ProductGroupItem))
}

type ProductGroupItem struct {
	ID           int64         `orm:"column(id);auto" json:"-"`
	Product      *Product      `orm:"column(product_id);null;rel(fk)" json:"product"`
	ProductGroup *ProductGroup `orm:"column(product_group_id);null;rel(fk)" json:"product_group"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ProductGroupItem) MarshalJSON() ([]byte, error) {
	type Alias ProductGroupItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ProductGroupItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ProductGroupItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
