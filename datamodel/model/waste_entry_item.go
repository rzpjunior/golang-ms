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
	orm.RegisterModel(new(WasteEntryItem))
}

// WasteEntryItem: struct to hold payment term model data for database
type WasteEntryItem struct {
	ID          int64   `orm:"column(id);auto" json:"-"`
	WasteQty    float64 `orm:"column(waste_qty)" json:"waste_qty"`
	Note        string  `orm:"column(note)" json:"note"`
	WasteReason int8    `orm:"column(waste_reason);null" json:"waste_reason,omitempty"`

	AvailableStock   float64 `orm:"-" json:"available_stock"`
	WasteStock       float64 `orm:"-" json:"waste_stock"`
	WasteReasonValue string  `orm:"-" json:"waste_reason_value"`

	WasteEntry *WasteEntry `orm:"column(waste_entry_id);null;rel(fk)" json:"waste_entry"`
	Product    *Product    `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *WasteEntryItem) MarshalJSON() ([]byte, error) {
	type Alias WasteEntryItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *WasteEntryItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *WasteEntryItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
