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
	orm.RegisterModel(new(PurchasePlanItem))
}

// PurchasePlanItem: struct to hold payment term model data for database
type PurchasePlanItem struct {
	ID              int64   `orm:"column(id);auto" json:"-"`
	PurchasePlanQty float64 `orm:"column(purchase_plan_qty)" json:"order_qty"`
	PurchaseQty     float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
	UnitPrice       float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal        float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight          float64 `orm:"column(weight)" json:"weight"`

	PurchasePlan       *PurchasePlan        `orm:"column(purchase_plan_id);null;rel(fk)" json:"purchase_plan"`
	Product            *Product             `orm:"column(product_id);null;rel(fk)" json:"product"`
	PurchaseOrderItems []*PurchaseOrderItem `orm:"reverse(many)" json:"purchase_order_items"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchasePlanItem) MarshalJSON() ([]byte, error) {
	type Alias PurchasePlanItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PurchasePlanItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchasePlanItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
