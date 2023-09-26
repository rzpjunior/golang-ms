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
	orm.RegisterModel(new(PurchaseOrderItem))
}

// PurchaseOrderItem: struct to hold payment term model data for database
type PurchaseOrderItem struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	OrderQty          float64 `orm:"column(order_qty)" json:"order_qty"`
	UnitPrice         float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal          float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight            float64 `orm:"column(weight)" json:"weight"`
	Note              string  `orm:"column(note)" json:"note"`
	MarketPurchaseStr string  `orm:"column(market_purchase)" json:"market_purchase_str"`
	PurchaseQty       float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
	TaxableItem       int8    `orm:"column(taxable_item)" json:"taxable_item"`
	IncludeTax        int8    `orm:"column(include_tax)" json:"include_tax"`
	TaxPercentage     float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	TaxAmount         float64 `orm:"column(tax_amount)" json:"tax_amount"`
	UnitPriceTax      float64 `orm:"column(unit_price_tax)" json:"unit_price_tax"`

	PurchaseOrder           *PurchaseOrder            `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order"`
	Product                 *Product                  `orm:"column(product_id);null;rel(fk)" json:"product"`
	FieldPurchaseOrderItems []*FieldPurchaseOrderItem `orm:"reverse(many)" json:"field_purchase_order_items"`
	PurchasePlanItem        *PurchasePlanItem         `orm:"column(purchase_plan_item_id);null;rel(fk)" json:"purchase_plan_item,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchaseOrderItem) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrderItem
	marketPurchase := []interface{}{}
	json.Unmarshal([]byte(m.MarketPurchaseStr), &marketPurchase)

	return json.Marshal(&struct {
		ID                 string        `json:"id"`
		IncludeTaxConvert  string        `json:"include_tax_convert"`
		TaxableItemConvert string        `json:"taxable_item_convert"`
		MarketPurchase     []interface{} `json:"market_purchase"`
		*Alias
	}{
		ID:                 common.Encrypt(m.ID),
		IncludeTaxConvert:  util.ConvertPurchaseOrderItemTaxStatus(m.IncludeTax),
		TaxableItemConvert: util.ConvertPurchaseOrderItemTaxStatus(m.TaxableItem),
		MarketPurchase:     marketPurchase,
		Alias:              (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PurchaseOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchaseOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
