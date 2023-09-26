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
	orm.RegisterModel(new(WasteDisposal))
}

// WasteDisposal: struct to hold model data for database
type WasteDisposal struct {
	ID              int64     `orm:"column(id);auto" json:"-"`
	Code            string    `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Note            string    `orm:"column(note)" json:"note"`
	Status          int8      `orm:"column(status)" json:"status"`

	// for print only
	City     string `orm:"-" json:"city"`
	Province string `orm:"-" json:"province"`

	Warehouse          *Warehouse           `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	WasteDisposalItems []*WasteDisposalItem `orm:"reverse(many)" json:"waste_disposal_item,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *WasteDisposal) MarshalJSON() ([]byte, error) {
	type Alias WasteDisposal

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *WasteDisposal) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *WasteDisposal) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
