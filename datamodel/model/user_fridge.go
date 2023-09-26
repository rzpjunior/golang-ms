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
	orm.RegisterModel(new(UserFridge))
}

// User model for user table.
type UserFridge struct {
	ID          int64      `orm:"column(id);auto" json:"-"`
	Code        string     `orm:"column(code);size(50);null" json:"code"`
	Username    string     `orm:"column(username);size(100);null" json:"username"`
	Password    string     `orm:"column(password);size(250);null" json:"password"`
	Token       string     `orm:"column(token);size(100);null" json:"token"`
	ForceLogout int8       `orm:"column(force_logout);null" json:"force_logout,omitempty"`
	Status      int8       `orm:"column(status);null" json:"status"`
	CreatedAt   time.Time  `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	LastLoginAt time.Time  `orm:"column(last_login_at);type(timestamp);null" json:"last_login_at"`
	Note        string     `orm:"column(note);size(250);null" json:"note"`
	Branch      *Branch    `orm:"column(branch_id);null;rel(fk)" json:"branch"`
	Warehouse   *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *UserFridge) MarshalJSON() ([]byte, error) {
	type Alias UserFridge

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// hide password user
	if m.Password != "" {
		m.Password = "********"
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *UserFridge) Save(fields ...string) (err error) {
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
func (m *UserFridge) Delete() (err error) {
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
func (m *UserFridge) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
