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
	orm.RegisterModel(new(CourierTransaction))
}

// CourierTransaction model for CourierTransaction table.
type CourierTransaction struct {
	ID             int64     `orm:"column(id);auto" json:"-"`
	Latitude       string    `orm:"column(latitude);size(250);null" json:"latitude"`
	Longitude      string    `orm:"column(longitude);size(250);null" json:"longitude"`
	Accuracy       string    `orm:"column(accuracy);size(250);null" json:"accuracy"`
	CourierName    string    `orm:"column(courier_name);size(100);null" json:"courier_name"`
	CourierPhoneNo string    `orm:"column(courier_phone_no);size(15);null" json:"courier_phone_no"`
	Note           string    `orm:"column(note);size(250);null" json:"note"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`

	DeliveryOrder *DeliveryOrder `orm:"column(delivery_order_id);null;rel(fk)" json:"delivery_order"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *CourierTransaction) MarshalJSON() ([]byte, error) {
	type Alias CourierTransaction

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *CourierTransaction) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *CourierTransaction) Delete() (err error) {
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
func (m *CourierTransaction) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
