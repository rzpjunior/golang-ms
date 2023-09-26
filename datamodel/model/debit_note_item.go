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
	orm.RegisterModel(new(DebitNoteItem))
}

// DebitNoteItem : struct to hold model data for database
type DebitNoteItem struct {
	ID          int64   `orm:"column(id);auto" json:"-"`
	ReturnQty   float64 `orm:"column(return_qty)" json:"return_qty"`
	ReceivedQty float64 `orm:"column(received_qty)" json:"received_qty"`
	UnitPrice   float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal    float64 `orm:"column(subtotal)" json:"subtotal"`
	Note        string  `orm:"column(note)" json:"note"`

	DebitNote *DebitNote `orm:"column(debit_note_id);null;rel(fk)" json:"debit_note,omitempty"`
	Product   *Product   `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DebitNoteItem) MarshalJSON() ([]byte, error) {
	type Alias DebitNoteItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DebitNoteItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DebitNoteItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
