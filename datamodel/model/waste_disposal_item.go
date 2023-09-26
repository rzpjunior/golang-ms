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
	orm.RegisterModel(new(WasteDisposalItem))
}

// WasteDisposalItem: struct to hold model data for database
type WasteDisposalItem struct {
	ID         int64   `orm:"column(id);auto" json:"-"`
	DisposeQty float64 `orm:"column(dispose_qty)" json:"dispose_qty"`
	Note       string  `orm:"column(note)" json:"note"`

	WasteDisposal *WasteDisposal `orm:"column(waste_disposal_id);null;rel(fk)" json:"waste_disposal"`
	Product       *Product       `orm:"column(product_id);null;rel(fk)" json:"product"`
	Stock         *Stock         `orm:"-" json:"stock"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *WasteDisposalItem) MarshalJSON() ([]byte, error) {
	type Alias WasteDisposalItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *WasteDisposalItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *WasteDisposalItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
