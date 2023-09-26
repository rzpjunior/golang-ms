// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(WasteLog))
}

// Waste Entry: struct to hold model data for database
type WasteLog struct {
	ID           int64   `orm:"column(id);auto" json:"-"`
	Ref          int64   `orm:"column(ref_id)" json:"ref"`
	RefType      int8    `orm:"column(ref_type)" json:"ref_type"`
	Type         int8    `orm:"column(type)" json:"type"`
	InitialStock float64 `orm:"column(initial_stock)" json:"initial_stock"`
	Quantity     float64 `orm:"column(quantity)" json:"quantity"`
	FinalStock   float64 `orm:"column(final_stock)" json:"final_stock"`
	UnitCost     float64 `orm:"column(unit_cost)" json:"unit_cost"`
	DocNote      string  `orm:"column(doc_note)" json:"doc_note"`
	ItemNote     string  `orm:"column(item_note)" json:"item_note"`
	Status       int8    `orm:"column(status)" json:"status"`
	WasteReason  int8    `orm:"column(waste_reason);null" json:"waste_reason,omitempty"`

	Warehouse *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	Product   *Product   `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *WasteLog) MarshalJSON() ([]byte, error) {
	type Alias WasteLog

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *WasteLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *WasteLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
