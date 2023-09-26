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
	orm.RegisterModel(new(SalesGroupItem))
}

// Sales Group Item: struct to hold model data for database
type SalesGroupItem struct {
	ID int64 `orm:"column(id);auto" json:"-"`

	SalesGroup  *SalesGroup  `orm:"column(sales_group_id);null;rel(fk)" json:"sales_group"`
	SubDistrict *SubDistrict `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesGroupItem) MarshalJSON() ([]byte, error) {
	type Alias SalesGroupItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesGroupItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesGroupItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
