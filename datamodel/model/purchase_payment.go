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
	orm.RegisterModel(new(PurchasePayment))
}

// PurchasePayment: struct to hold model data for database
type PurchasePayment struct {
	ID                       int64     `orm:"column(id);auto" json:"-"`
	Code                     string    `orm:"column(code);size(50);null" json:"code"`
	Status                   int8      `orm:"column(status);null" json:"status"`
	RecognitionDate          time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Amount                   float64   `orm:"column(amount)" json:"amount"`
	PaidOff                  int8      `orm:"column(paid_off);null" json:"paid_off"`
	Note                     string    `orm:"column(note)" json:"note"`
	ImageUrl                 string    `orm:"column(image_url)" json:"image_url"`
	BankPaymentVoucherNumber string    `orm:"column(bank_payment_voucher_number)" json:"bank_payment_voucher_number"`

	PurchaseInvoice *PurchaseInvoice `orm:"column(purchase_invoice_id);null;rel(fk)" json:"purchase_invoice"`
	PaymentMethod   *PaymentMethod   `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method"`

	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchasePayment) MarshalJSON() ([]byte, error) {
	type Alias PurchasePayment

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PurchasePayment) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchasePayment) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
