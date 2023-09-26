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
	orm.RegisterModel(new(SalesInvoiceItem))
}

// Sales Invoice Item: struct to hold model sales_invoice_item for database
type SalesInvoiceItem struct {
	ID            int64   `orm:"column(id);auto" json:"-"`
	InvoiceQty    float64 `orm:"column(invoice_qty)" json:"invoice_qty"`
	UnitPrice     float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal      float64 `orm:"column(subtotal)" json:"subtotal"`
	Note          string  `orm:"column(note)" json:"note"`
	TaxableItem   int8    `orm:"column(taxable_item)" json:"taxable_item"`
	TaxPercentage float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	SkuDiscAmount float64 `orm:"column(sku_disc_amount)" json:"sku_disc_amount"`

	SalesInvoice   *SalesInvoice   `orm:"column(sales_invoice_id);null;rel(fk)" json:"sales_invoice"`
	SalesOrderItem *SalesOrderItem `orm:"column(sales_order_item_id);null;rel(fk)" json:"sales_order_item"`
	Product        *Product        `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesInvoiceItem) MarshalJSON() ([]byte, error) {
	type Alias SalesInvoiceItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesInvoiceItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesInvoiceItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *SalesInvoiceItem) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}
