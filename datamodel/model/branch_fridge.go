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
	orm.RegisterModel(new(BranchFridge))
}

// BranchFridge model for BranchFridge table.
type BranchFridge struct {
	ID                    int64      `orm:"column(id);auto" json:"-"`
	Code                  string     `orm:"column(code);size(50);null" json:"code"`
	Note                  string     `orm:"column(note);size(250);null" json:"note"`
	Status                int8       `orm:"column(status);null" json:"status"`
	MacAddress            string     `orm:"column(mac_address);null" json:"mac_address"`
	CreatedAt             time.Time  `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy             int64      `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt         time.Time  `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy         int64      `orm:"column(last_updated_by)" json:"last_updated_by"`
	LastSeenAt            time.Time  `orm:"column(last_seen_at);type(timestamp);null" json:"last_seen_at"`
	Branch                *Branch    `orm:"column(branch_id);null;rel(fk)"`
	Warehouse             *Warehouse `orm:"column(warehouse_id);null;rel(fk)" `
	Reader                string     `orm:"-" json:"reader"`
	ReplenishmentRequired string     `orm:"-" json:"replenishment_required"`

	FridgeProductDemand *FridgeProductDemand `orm:"-" json:"fridge_product_demand,omitempty"`
}

type BranchFridgeListQuery struct {
	CustomerName          string          `orm:"column(customer_name);null" json:"customer_name,omitempty"`
	LastSeenAt            time.Time       `orm:"column(last_seen_at);null"  json:"last_seen_at,omitempty"`
	Status                string          `orm:"column(status);null"  json:"status,omitempty"`
	WarehouseName         string          `orm:"column(warehouse_name);null" json:"warehouse_name,omitempty"`
	BranchName            string          `orm:"column(branch_name);null" json:"branch_name,omitempty"`
	AllFridge             []*BranchFridge `orm:"column(all_fridge);null" json:"all_fridge,omitempty"`
	ReplenishmentRequired string          `orm:"-" json:"replenishment_required"`
	Products              int             `orm:"-" json:"products"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *BranchFridge) MarshalJSON() ([]byte, error) {
	type Alias BranchFridge

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating Branch Fridge struct into Branch Fridge table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Branch Fridge.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *BranchFridge) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Branch Fridge data
// this also will truncated all data from all table
// that have relation with this Branch Fridge.
func (m *BranchFridge) Delete() (err error) {
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
func (m *BranchFridge) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
