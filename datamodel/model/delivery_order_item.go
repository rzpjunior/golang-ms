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
	orm.RegisterModel(new(DeliveryOrderItem))
}

// DeliveryOrderItem: struct to hold model data for database
type DeliveryOrderItem struct {
	ID              int64   `orm:"column(id);auto" json:"-"`
	DeliverQty      float64 `orm:"column(deliver_qty)" json:"deliver_qty"`
	ReceiveQty      float64 `orm:"column(receive_qty)" json:"receive_qty"`
	ReceiptItemNote string  `orm:"column(receipt_item_note)" json:"receipt_item_note"`
	OrderItemNote   string  `orm:"column(order_item_note)" json:"order_item_note"`
	Weight          float64 `orm:"column(weight)" json:"weight"`
	Note            string  `orm:"column(note)" json:"note"`

	OldQuantity float64 `orm:"-" json:"-"`

	DeliveryOrder  *DeliveryOrder  `orm:"column(delivery_order_id);null;rel(fk)" json:"delivery_order"`
	SalesOrderItem *SalesOrderItem `orm:"column(sales_order_item_id);null;rel(fk)" json:"sales_order_item"`
	Product        *Product        `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DeliveryOrderItem) MarshalJSON() ([]byte, error) {
	type Alias DeliveryOrderItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Delete permanently deleting item data
// this also will truncated all data from all table
// that have relation with this item.
func (m *DeliveryOrderItem) Delete() (err error) {
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

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DeliveryOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DeliveryOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
