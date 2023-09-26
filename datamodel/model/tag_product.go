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
	orm.RegisterModel(new(TagProduct))
}

// TagProduct : struct to hold model data for database
type TagProduct struct {
	ID       int64  `orm:"column(id);auto" json:"-"`
	Area     string `orm:"column(area)" json:"area"`
	Code     string `orm:"column(code)" json:"code"`
	Name     string `orm:"column(name)" json:"name"`
	Value    string `orm:"column(value)" json:"value"`
	ImageUrl string `orm:"column(image_url)" json:"image_url"`
	Note     string `orm:"column(note)" json:"note"`
	Status   int8   `orm:"column(status)" json:"status"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *TagProduct) MarshalJSON() ([]byte, error) {
	type Alias TagProduct

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *TagProduct) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *TagProduct) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
