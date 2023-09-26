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
	orm.RegisterModel(new(FridgeReplenishmentItem))
}

// FridgeReplenishmentItem model for Fridge Replenishment Item table.
type FridgeReplenishmentItem struct {
	ID                  int64                `orm:"column(id);auto" json:"-"`
	FridgeReplenishment *FridgeReplenishment `orm:"column(fridge_replenishment_id);null;rel(fk)" json:"fridge_replenishment,omitempty"`
	Product             *Product             `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
	RequestedQty        float64              `orm:"column(requested_qty)" json:"requested_qty"`
	ProductDemandQty    float64              `orm:"-" json:"product_demand_qty"`
	RemainingQty        float64              `orm:"column(remaining_qty)" json:"remaining_qty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *FridgeReplenishmentItem) MarshalJSON() ([]byte, error) {
	type Alias FridgeReplenishmentItem

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating Fridge struct into Fridge table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Fridge.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *FridgeReplenishmentItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Fridge data
// this also will truncated all data from all table
// that have relation with this Fridge.
func (m *FridgeReplenishmentItem) Delete() (err error) {
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

// Read execute select based on data struct that already
// assigned.
func (m *FridgeReplenishmentItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
