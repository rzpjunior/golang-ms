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
	orm.RegisterModel(new(TxnXendit))
}

// TxnXendit model for Settlement table.
type TxnXendit struct {
	ID              int64           `orm:"column(id);auto" json:"-"`
	Merchant        *Merchant       `orm:"column(merchant_id);null;rel(fk)" json:"merchant,omitempty"`
	SalesOrder      *SalesOrder     `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order,omitempty"`
	PaymentChannel  *PaymentChannel `orm:"column(payment_channel_id);null;rel(fk)" json:"payment_channel,omitempty"`
	Type            int             `orm:"column(type);null" json:"type"`
	AccountNumber   string          `orm:"column(account_number);null" json:"account_number"`
	Amount          float64         `orm:"column(amount);null" json:"amount"`
	TransactionDate string          `orm:"column(transaction_date);null" json:"transaction_date"`
	TransactionTime string          `orm:"column(transaction_time);null" json:"transaction_time"`
	CreatedAt       time.Time       `orm:"column(created_at);type(timestamp);null" json:"created_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *TxnXendit) MarshalJSON() ([]byte, error) {
	type Alias TxnXendit

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating Promotion struct into promotion table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to promotion.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *TxnXendit) Save(fields ...string) (err error) {
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
func (m *TxnXendit) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
