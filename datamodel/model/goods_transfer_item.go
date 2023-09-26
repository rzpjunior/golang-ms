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
	orm.RegisterModel(new(GoodsTransferItem))
}

// GoodsTransferItem: struct to hold model data for database
type GoodsTransferItem struct {
	ID          int64   `orm:"column(id);auto" json:"-"`
	DeliverQty  float64 `orm:"column(deliver_qty)" json:"delivery_qty"`
	ReceiveQty  float64 `orm:"column(receive_qty)" json:"receive_qty"`
	RequestQty  float64 `orm:"column(request_qty)" json:"request_qty"`
	ReceiveNote string  `orm:"column(receive_note)" json:"receive_note"`
	UnitCost    float64 `orm:"column(unit_cost)" json:"unit_cost"`
	Subtotal    float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight      float64 `orm:"column(weight)" json:"weight"`
	Note        string  `orm:"column(note)" json:"note"`

	GoodsTransfer *GoodsTransfer `orm:"column(goods_transfer_id);null;rel(fk)" json:"goods_transfer"`
	Product       *Product       `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *GoodsTransferItem) MarshalJSON() ([]byte, error) {
	type Alias GoodsTransferItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *GoodsTransferItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *GoodsTransferItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
