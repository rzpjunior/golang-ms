// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(DashboardArea))
}

type DashboardArea struct {
	ID   int64  `orm:"column(id)" json:"id"`
	Name string `orm:"column(name)" json:"name"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DashboardArea) MarshalJSON() ([]byte, error) {
	type Alias DashboardArea

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DashboardArea) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	o.Using("scrape")

	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DashboardArea) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("scrape")

	return o.Read(m, fields...)
}
