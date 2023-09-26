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
	orm.RegisterModel(new(Fridge))
}

// Fridge model for Fridge table.
type Fridge struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	Code          string    `orm:"column(code);size(50);null" json:"code"`
	Name          string    `orm:"column(name);size(100);null" json:"name"`
	Note          string    `orm:"column(note);size(250);null" json:"note"`
	Status        int8      `orm:"column(status);null" json:"status"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Fridge) MarshalJSON() ([]byte, error) {
	type Alias Fridge

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
func (m *Fridge) Save(fields ...string) (err error) {
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
func (m *Fridge) Delete() (err error) {
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
func (m *Fridge) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
