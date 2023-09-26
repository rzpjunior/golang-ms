// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(FridgeReplenishment))
}

// FridgeReplenishment model for Fridge Replenishment table.
type FridgeReplenishment struct {
	ID        int64      `orm:"column(id);auto" json:"-"`
	Code      string     `orm:"column(code);size(50);null" json:"code"`
	Warehouse *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	Branch    *Branch    `orm:"-" json:"branch,omitempty"`

	Wrt          *Wrt      `orm:"column(wrt_id);null;rel(fk)" json:"wrt,omitempty"`
	DeliveryDate time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	Note         string    `orm:"column(note);size(250);null" json:"note"`
	Status       int8      `orm:"column(status)" json:"status"`

	FridgeReplenishmentItems []*FridgeReplenishmentItem `orm:"reverse(many)" json:"fridge_replenishment_items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *FridgeReplenishment) MarshalJSON() ([]byte, error) {
	type Alias FridgeReplenishment

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
func (m *FridgeReplenishment) Save(fields ...string) (err error) {
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
func (m *FridgeReplenishment) Delete() (err error) {
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
func (m *FridgeReplenishment) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
