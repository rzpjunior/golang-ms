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
	orm.RegisterModel(new(SupplierReturn))
}

// SupplierReturn : struct to hold model data for database
type SupplierReturn struct {
	ID              int64     `orm:"column(id);auto" json:"-"`
	Code            string    `orm:"column(code)" json:"code"`
	RecognitionDate time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Note            string    `orm:"column(note)" json:"note"`
	Status          int       `orm:"column(status)" json:"status"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy       *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	ConfirmedAt     time.Time `orm:"column(confirmed_at);type(timestamp);null" json:"confirmed_at"`
	ConfirmedBy     *Staff    `orm:"column(confirmed_by);null;rel(fk)" json:"confirmed_by"`
	DeltaPrint      int8      `orm:"column(delta_print)" json:"delta_print"`
	ReturnType      int8      `orm:"column(return_type)" json:"return_type"`

	Supplier     *Supplier     `orm:"column(supplier_id);null;rel(fk)" json:"supplier,omitempty"`
	GoodsReceipt *GoodsReceipt `orm:"column(goods_receipt_id);null;rel(fk)" json:"good_receipt,omitempty"`
	Warehouse    *Warehouse    `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`

	SupplierReturnItems []*SupplierReturnItem `orm:"reverse(many)" json:"supplier_return_items,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SupplierReturn) MarshalJSON() ([]byte, error) {
	type Alias SupplierReturn

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SupplierReturn) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SupplierReturn) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
