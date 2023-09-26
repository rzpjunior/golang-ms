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
	orm.RegisterModel(new(PurchaseInvoice))
}

// PurchaseInvoice: struct to hold model data for database
type PurchaseInvoice struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	Status           int8      `orm:"column(status);null" json:"status"`
	RecognitionDate  time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	DueDate          time.Time `orm:"column(due_date)" json:"due_date"`
	TaxPct           float64   `orm:"column(tax_pct)" json:"tax_pct"`
	DeliveryFee      float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	Adjustment       int8      `orm:"column(adjustment);null" json:"adjustment"`
	AdjAmount        float64   `orm:"column(adj_amount)" json:"adj_amount"`
	AdjNote          string    `orm:"column(adj_note)" json:"adj_note"`
	TotalPrice       float64   `orm:"column(total_price)" json:"total_price"`
	TaxAmount        float64   `orm:"column(tax_amount)" json:"tax_amount"`
	TotalCharge      float64   `orm:"column(total_charge)" json:"total_charge"`
	Note             string    `orm:"column(note)" json:"note"`
	TaxInvoiceURL    string    `orm:"column(tax_invoice_url)" json:"tax_invoice_url"`
	TaxInvoiceNumber string    `orm:"column(tax_invoice_number)" json:"tax_invoice_number"`
	RemainingAmount  float64   `orm:"-" json:"remaining_amount"`
	IsPaid           int8      `orm:"-" json:"is_paid"`
	DebitNoteIDs     string    `orm:"column(debit_note_id)" json:"debit_note_id"`

	PurchaseOrder *PurchaseOrder `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order"`
	PurchaseTerm  *PurchaseTerm  `orm:"column(term_payment_pur_id);null;rel(fk)" json:"purchase_term"`
	DebitNote     []*DebitNote   `orm:"-" json:"debit_note"`

	PurchaseInvoiceItems []*PurchaseInvoiceItem `orm:"reverse(many)" json:"purchase_invoice_items,omitempty"`

	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchaseInvoice) MarshalJSON() ([]byte, error) {
	type Alias PurchaseInvoice

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PurchaseInvoice) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchaseInvoice) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
