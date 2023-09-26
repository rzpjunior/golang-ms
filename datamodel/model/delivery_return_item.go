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
	orm.RegisterModel(new(DeliveryReturnItem))
}

// DeliveryReturnItem: struct to hold model data for database
type DeliveryReturnItem struct {
	ID               int64   `orm:"column(id);auto" json:"-"`
	ReturnGoodQty    float64 `orm:"column(return_good_qty)" json:"return_good_qty"`
	ReturnWasteQty   float64 `orm:"column(return_waste_qty)" json:"return_waste_qty"`
	WasteReason      int8    `orm:"column(waste_reason);null" json:"return_waste_reason"`
	WasteReasonValue string  `orm:"-" json:"return_waste_reason_value"`
	UnitCost         float64 `orm:"column(unit_cost)" json:"unit_cost"`
	Note             string  `orm:"column(note)" json:"note"`

	DeliveryReturn    *DeliveryReturn    `orm:"column(delivery_return_id);null;rel(fk)" json:"delivery_return"`
	DeliveryOrderItem *DeliveryOrderItem `orm:"column(delivery_order_item_id);null;rel(fk)" json:"delivery_order_item"`
	Product           *Product           `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DeliveryReturnItem) MarshalJSON() ([]byte, error) {
	type Alias DeliveryReturnItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DeliveryReturnItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DeliveryReturnItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
