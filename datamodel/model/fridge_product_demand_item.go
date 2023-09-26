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
	orm.RegisterModel(new(FridgeProductDemandItem))
}

// FridgeProductDemandItem model for Fridge Product Demand Item table.
type FridgeProductDemandItem struct {
	ID                    int64                `orm:"column(id);auto" json:"-"`
	FridgeProductDemand   *FridgeProductDemand `orm:"column(fridge_product_demand_id);null;rel(fk)" json:"fridge_product_demand,omitempty"`
	Product               *Product             `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
	RequestedQty          float64              `orm:"column(requested_qty)" json:"requested_qty"`
	ThresholdQty          float64              `orm:"column(threshold_qty)" json:"threshold_qty"`
	RemainingQty          float64              `orm:"-" json:"remaining_qty"`
	ReplenishmentRequired bool                 `orm:"-" json:"replenishment_required"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *FridgeProductDemandItem) MarshalJSON() ([]byte, error) {
	type Alias FridgeProductDemandItem

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
func (m *FridgeProductDemandItem) Save(fields ...string) (err error) {
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
func (m *FridgeProductDemandItem) Delete() (err error) {
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
func (m *FridgeProductDemandItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
