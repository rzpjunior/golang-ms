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
	orm.RegisterModel(new(StockOpname))
}

// StockOpname: struct to hold model data for database
type StockOpname struct {
	ID               int64              `orm:"column(id);auto" json:"-"`
	Code             string             `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate  time.Time          `orm:"column(recognition_date)" json:"recognition_date"`
	Note             string             `orm:"column(note)" json:"note"`
	Status           int8               `orm:"column(status);null" json:"status"`
	StockType        int8               `orm:"column(stock_type)" json:"stock_type"`
	CancellationNote string             `orm:"-" json:"cancellation_note,omitempty"`
	StockOpnameItems []*StockOpnameItem `orm:"reverse(many)" json:"stock_opname_items,omitempty"`

	Category  *Category  `orm:"column(category_id);null;rel(fk)" json:"category"`
	Warehouse *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *StockOpname) MarshalJSON() ([]byte, error) {
	type Alias StockOpname

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *StockOpname) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *StockOpname) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
