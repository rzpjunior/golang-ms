// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
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
	orm.RegisterModel(new(DebitNote))
}

// DebitNote : struct to hold model data for database
type DebitNote struct {
	ID                    int64     `orm:"column(id);auto" json:"-"`
	Code                  string    `orm:"column(code)" json:"code"`
	RecognitionDate       time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	TotalPrice            float64   `orm:"column(total_price)" json:"total_price"`
	Note                  string    `orm:"column(note)" json:"note"`
	Status                int8      `orm:"column(status)" json:"status"`
	UsedInPurchaseInvoice int8      `orm:"column(used_in_purchase_invoice)" json:"used_in_purchase_invoice"`
	CreatedAt             time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy             *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`

	SupplierReturn *SupplierReturn `orm:"column(supplier_return_id);null;rel(fk)" json:"supplier_return,omitempty"`

	DebitNoteItems []*DebitNoteItem `orm:"reverse(many)" json:"debit_note_items,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DebitNote) MarshalJSON() ([]byte, error) {
	type Alias DebitNote

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DebitNote) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DebitNote) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
