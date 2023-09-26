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
	orm.RegisterModel(new(ConfigApp))
}

// ConfigApp model for config_app table.
type ConfigApp struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Application int64  `orm:"column(application);size(45)" json:"application,omitempty"`
	Field       string `orm:"column(field);null" json:"field,omitempty"`
	Attribute   string `orm:"column(attribute);null" json:"attribute,omitempty"`
	Value       string `orm:"column(value);null" json:"value"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ConfigApp) MarshalJSON() ([]byte, error) {
	type Alias ConfigApp

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating ConfigApp struct into app_config table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to app_config.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ConfigApp) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting app_config data
// this also will truncated all data from all table
// that have relation with this app_config.
func (m *ConfigApp) Delete() (err error) {
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
func (m *ConfigApp) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
