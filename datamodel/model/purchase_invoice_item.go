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
	orm.RegisterModel(new(PurchaseInvoiceItem))
}

// PurchaseInvoiceItem: struct to hold model data for database
type PurchaseInvoiceItem struct {
	ID            int64   `orm:"column(id);auto" json:"-"`
	InvoiceQty    float64 `orm:"column(invoice_qty)" json:"invoice_qty"`
	UnitPrice     float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal      float64 `orm:"column(subtotal)" json:"subtotal"`
	Note          string  `orm:"column(note)" json:"note"`
	TaxableItem   int8    `orm:"column(taxable_item)" json:"taxable_item"`
	IncludeTax    int8    `orm:"column(include_tax)" json:"include_tax"`
	TaxPercentage float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	TaxAmount     float64 `orm:"column(tax_amount)" json:"tax_amount"`
	UnitPriceTax  float64 `orm:"column(unit_price_tax)" json:"unit_price_tax"`

	PurchaseInvoice   *PurchaseInvoice   `orm:"column(purchase_invoice_id);null;rel(fk)" json:"purchase_invoice"`
	PurchaseOrderItem *PurchaseOrderItem `orm:"column(purchase_order_item_id);null;rel(fk)" json:"purchase_order_item"`
	Product           *Product           `orm:"column(product_id);null;rel(fk)" json:"product"`

	GoodsReceiptItem *GoodsReceiptItem `orm:"-"json:"good_receipt_item"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchaseInvoiceItem) MarshalJSON() ([]byte, error) {
	type Alias PurchaseInvoiceItem

	return json.Marshal(&struct {
		ID                 string `json:"id"`
		IncludeTaxConvert  string `json:"include_tax_convert"`
		TaxableItemConvert string `json:"taxable_item_convert"`
		*Alias
	}{
		ID:                 common.Encrypt(m.ID),
		IncludeTaxConvert:  util.ConvertPurchaseInvoiceItemTaxStatus(m.IncludeTax),
		TaxableItemConvert: util.ConvertPurchaseInvoiceItemTaxStatus(m.TaxableItem),
		Alias:              (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PurchaseInvoiceItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchaseInvoiceItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
