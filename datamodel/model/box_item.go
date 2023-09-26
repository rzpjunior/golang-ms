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
	orm.RegisterModel(new(BoxItem))
}

// BoxItem model for BoxItem table.
type BoxItem struct {
	ID int64 `orm:"column(id);auto" json:"-"`
	//Code          string    `orm:"column(code);size(50);null" json:"code"`
	Note          string    `orm:"column(note);size(250);null" json:"note"`
	UnitPrice     float64   `orm:"column(unit_price)" json:"unit_price"`
	TotalPrice    float64   `orm:"column(total_price)" json:"total_price"`
	TotalWeight   float64   `orm:"column(total_weight)" json:"total_weight"`
	Status        int8      `orm:"column(status);null" json:"status"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
	//LastSeenAt    time.Time `orm:"column(last_seen_at);type(timestamp);null" json:"last_seen_at"`
	FinishedAt time.Time `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
	FinishedBy int64     `orm:"column(finished_by)" json:"finished_by"`

	Box     *Box     `orm:"column(box_id);null;rel(fk)" json:"box,omitempty"`
	Product *Product `orm:"column(product_id);null;rel(fk)" json:"product"`
}

type ProductFridgeBoxListQuery struct {
	ProductName     string    `orm:"column(product_name);null" json:"product_name,omitempty"`
	TotalWeight     float64   `orm:"column(total_weight);null" json:"total_weight,omitempty"`
	ItemImage       string    `orm:"column(item_image);null" json:"item_image,omitempty"`
	Uom             string    `orm:"column(uom_name);null" json:"uom_name,omitempty"`
	ProcessedAt     time.Time `orm:"column(processed_at);null"  json:"processed_at,omitempty"`
	Rfid            string    `orm:"column(rfid);null"  json:"rfid,omitempty"`
	WarehouseId     int64     `orm:"column(warehouse_id);null"  json:"-"`
	WasteImage      string    `orm:"column(waste_image);null" json:"waste_image,omitempty"`
	FinishedAt      time.Time `orm:"column(finished_at);null"  json:"finished_at,omitempty"`
	BoxFridgeStatus int64     `orm:"column(box_fridge_status);null"  json:"box_fridge_status,omitempty"`
	BoxItemStatus   int64     `orm:"column(box_item_status);null"  json:"box_item_status,omitempty"`
	Status          string    `orm:"column(status);null"  json:"status,omitempty"`
	WarehouseName   string    `orm:"column(warehouse_name);null" json:"warehouse_name,omitempty"`
	BranchName      string    `orm:"column(branch_name);null" json:"branch_name,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *BoxItem) MarshalJSON() ([]byte, error) {
	type Alias BoxItem

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating Box Item struct into Box Item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Box Item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *BoxItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Box Item data
// this also will truncated all data from all table
// that have relation with this Box Item.
func (m *BoxItem) Delete() (err error) {
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
func (m *BoxItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
