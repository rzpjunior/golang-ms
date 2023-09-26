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
	orm.RegisterModel(new(BoxFridge))
}

// BoxFridge model for BoxFridge table.
type BoxFridge struct {
	ID int64 `orm:"column(id);auto" json:"-"`
	//Code          string     `orm:"column(code);size(50);null" json:"code"`
	Note          string     `orm:"column(note);size(250);null" json:"note"`
	Status        int8       `orm:"column(status);null" json:"status"`
	ImageUrl      string     `orm:"column(image_url);size(300);null" json:"image_url"`
	CreatedAt     time.Time  `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     int64      `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt time.Time  `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastSeenAt    time.Time  `orm:"column(last_seen_at);type(timestamp);null" json:"last_seen_at"`
	LastUpdatedBy int64      `orm:"column(last_updated_by)" json:"last_updated_by"`
	Box           *Box       `orm:"column(box_id);null;rel(fk)"`
	Warehouse     *Warehouse `orm:"column(warehouse_id);null;rel(fk)" `
	BoxItem       *BoxItem   `orm:"column(box_item_id);null;rel(fk)"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *BoxFridge) MarshalJSON() ([]byte, error) {
	type Alias BoxFridge

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating Box Fridge struct into Box Fridge table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Box Fridge.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *BoxFridge) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Box Fridge data
// this also will truncated all data from all table
// that have relation with this Box Fridge.
func (m *BoxFridge) Delete() (err error) {
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
func (m *BoxFridge) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
