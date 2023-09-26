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
	orm.RegisterModel(new(StockOpnameItem))
}

// StockOpnameItem: struct to hold payment term model data for database
type StockOpnameItem struct {
	ID           int64   `orm:"column(id);auto" json:"-"`
	InitialStock float64 `orm:"column(initial_stock)" json:"initial_stock"`
	AdjustQty    float64 `orm:"column(adjust_qty)" json:"adjust_qty"`
	FinalStock   float64 `orm:"column(final_stock)" json:"final_stock"`
	OpnameReason int8    `orm:"column(opname_reason)" json:"opname_reason"`
	Note         string  `orm:"column(note)" json:"note"`

	OpnameReasonValue string `orm:"-" json:"opname_reason_value"`

	StockOpname *StockOpname `orm:"column(stock_opname_id);null;rel(fk)" json:"stock_opname"`
	Product     *Product     `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *StockOpnameItem) MarshalJSON() ([]byte, error) {
	type Alias StockOpnameItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *StockOpnameItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *StockOpnameItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
