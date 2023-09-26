// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
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
	orm.RegisterModel(new(SalesInvoiceExternal))
}

// SalesInvoiceExternal model for district table.
type SalesInvoiceExternal struct {
	ID              int64       `orm:"column(id);auto" json:"-"`
	SalesOrder      *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order,omitempty"`
	XenditInvoiceID string      `orm:"column(xendit_invoice_id);size(35);null" json:"xendit_invoice_id"`
	CreatedAt       time.Time   `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CancelledAt     time.Time   `orm:"column(cancelled_at);type(timestamp);null" json:"cancelled_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesInvoiceExternal) MarshalJSON() ([]byte, error) {
	type Alias SalesInvoiceExternal

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SalesInvoiceExternal) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read execute select based on data struct that already
// assigned.
func (m *SalesInvoiceExternal) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
