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
	orm.RegisterModel(new(SalesOrderItem))
}

// Sales Order Item: struct to hold model sales_order_item for database
type SalesOrderItem struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	OrderQty          float64 `orm:"column(order_qty)" json:"order_qty"`
	UnitPrice         float64 `orm:"column(unit_price)" json:"unit_price"`
	ShadowPrice       float64 `orm:"column(shadow_price)" json:"shadow_price"`
	Subtotal          float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight            float64 `orm:"column(weight)" json:"weight"`
	Note              string  `orm:"column(note)" json:"note"`
	ProductPush       int8    `orm:"column(product_push)" json:"product_push"`
	TaxableItem       int8    `orm:"column(taxable_item)" json:"taxable_item"`
	TaxPercentage     float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	DiscountQty       float64 `orm:"column(discount_qty);digits(10);decimals(2)" json:"discount_qty"`
	UnitPriceDiscount float64 `orm:"column(unit_price_discount);digits(10);decimals(2)" json:"unit_price_discount"`
	SkuDiscountAmount float64 `orm:"column(sku_disc_amount)" json:"sku_disc_amount"`
	DefaultPrice      float64 `orm:"column(default_price)" json:"default_price"`

	SalesOrder        *SalesOrder        `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
	Product           *Product           `orm:"column(product_id);null;rel(fk)" json:"product"`
	SkuDiscountItem   *SkuDiscountItem   `orm:"column(sku_discount_item_id);null;rel(fk)" json:"sku_discount_item,omitempty"`
	SalesInvoiceItem  *SalesInvoiceItem  `orm:"-" json:"sales_invoice_item,omitempty"`
	DeliveryOrderItem *DeliveryOrderItem `orm:"-" json:"delivery_order_item,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesOrderItem) MarshalJSON() ([]byte, error) {
	type Alias SalesOrderItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
