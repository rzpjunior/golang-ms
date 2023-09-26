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
	orm.RegisterModel(new(Box))
}

// Box model for Box table.
type Box struct {
	ID   int64  `orm:"column(id);auto" json:"-"`
	Rfid string `orm:"column(rfid);size(50);null" json:"rfid"`
	//Code string `orm:"column(code);size(50);null" json:"code"`
	//Name          string    `orm:"column(name);size(100);null" json:"name"`
	Note          string    `orm:"column(note);size(250);null" json:"note"`
	Status        int8      `orm:"column(status);null" json:"status"`
	Size          int8      `orm:"column(size);null" json:"size"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Box) MarshalJSON() ([]byte, error) {
	type Alias Box

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating Box struct into Box table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Box.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Box) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Box data
// this also will truncated all data from all table
// that have relation with this Box.
func (m *Box) Delete() (err error) {
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
func (m *Box) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
