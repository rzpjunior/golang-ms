// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(PublicArea2))
}

type PublicArea2 struct {
	ID   int64  `orm:"column(id)" json:"id"`
	Code string `orm:"column(code)" json:"code"`
	Name string `orm:"column(name)" json:"name"`
}

// TableName : set table name used by model
func (PublicArea2) TableName() string {
	return "public_area_2"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PublicArea2) MarshalJSON() ([]byte, error) {
	type Alias PublicArea2

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PublicArea2) Save(fields ...string) (err error) {
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
func (m *PublicArea2) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("scrape")

	return o.Read(m, fields...)
}
