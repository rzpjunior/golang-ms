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
	orm.RegisterModel(new(PickingOrderItem))
}

// PickingOrderItem model for picking order item table.
type PickingOrderItem struct {
	ID             int64   `orm:"column(id);auto" json:"-"`
	OrderQuantity  float64 `orm:"column(order_qty)" json:"order_qty"`
	PickQuantity   float64 `orm:"column(pick_qty)" json:"pick_qty"`
	CheckQuantity  float64 `orm:"column(check_qty)" json:"check_qty"`
	UnfullfillNote string  `orm:"column(unfullfill_note)" json:"unfullfill_note"`
	FlagOrder      int8    `orm:"column(flag_order)" json:"flag_order"`
	FlagSavePick   int8    `orm:"column(flag_saved_pick)" json:"flag_saved_pick"`
	PickingFlag    int8    `orm:"column(picking_flag)" json:"picking_flag"`

	SalesOrderItemNote string `orm:"-" json:"sales_order_item_note"`

	PickingOrderAssign *PickingOrderAssign `orm:"column(picking_order_assign_id);null;rel(fk)" json:"picking_order_assign,omitempty"`
	Product            *Product            `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PickingOrderItem) MarshalJSON() ([]byte, error) {
	type Alias PickingOrderItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PickingOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PickingOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
