// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(Permission))
}

// Permission model for permission table.
type Permission struct {
	ID            int64       `orm:"column(id);auto" json:"-"`
	Parent        *Permission `orm:"column(parent_id);null;rel(fk)" json:"parent,omitempty"`
	Code          string      `orm:"column(code);size(50);null" json:"code"`
	Value         string      `orm:"column(value);size(50);null" json:"value"`
	Name          string      `orm:"column(name);size(100);null" json:"name"`
	Note          string      `orm:"column(note);size(250);null" json:"note"`
	Status        int8        `orm:"column(status);null" json:"status"`
	StatusConvert string      `orm:"-" json:"status_convert"`

	Child []*Child `orm:"-" json:"child"`
}

// Child model for permission table.
type Child struct {
	ID         string        `orm:"column(id)" json:"id"`
	Parent     *Permission   `orm:"column(parent_id);null;rel(fk)" json:"parent,omitempty"`
	Code       string        `orm:"column(code);size(50);null" json:"code"`
	Value      string        `orm:"column(value);size(50);null" json:"value"`
	Name       string        `orm:"column(name);size(100);null" json:"name"`
	Note       string        `orm:"column(note);size(250);null" json:"note"`
	Status     int8          `orm:"column(status);null" json:"status"`
	GrandChild []*GrandChild `orm:"-" json:"grand_child"`
}

// GrandChild model for permission table.
type GrandChild struct {
	ID     string      `orm:"column(id)" json:"id"`
	Parent *Permission `orm:"column(parent_id);null;rel(fk)" json:"parent,omitempty"`
	Code   string      `orm:"column(code);size(50);null" json:"code"`
	Value  string      `orm:"column(value);size(50);null" json:"value"`
	Name   string      `orm:"column(name);size(100);null" json:"name"`
	Note   string      `orm:"column(note);size(250);null" json:"note"`
	Status int8        `orm:"column(status);null" json:"status"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Permission) MarshalJSON() ([]byte, error) {
	type Alias Permission

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Permission) Save(fields ...string) (err error) {
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
func (m *Permission) Delete() (err error) {
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
func (m *Permission) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
