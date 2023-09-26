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
	orm.RegisterModel(new(SupplierReturnItem))
}

// SupplierReturn : struct to hold model data for database
type SupplierReturnItem struct {
	ID             int64   `orm:"column(id);auto" json:"-"`
	ReturnGoodsQty float64 `orm:"column(return_good_qty)" json:"return_good_qty"`
	ReceivedQty    float64 `orm:"column(received_qty)" json:"received_qty"`
	Note           string  `orm:"column(note)" json:"note"`

	SupplierReturn *SupplierReturn `orm:"column(supplier_return_id);null;rel(fk)" json:"supplier_return,omitempty"`
	Product        *Product        `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SupplierReturnItem) MarshalJSON() ([]byte, error) {
	type Alias SupplierReturnItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SupplierReturnItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SupplierReturnItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
