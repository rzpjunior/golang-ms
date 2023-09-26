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
	orm.RegisterModel(new(BusinessType))
}

// BusinessType : struct to hold model data for database
type BusinessType struct {
	ID           int64  `orm:"column(id);auto" json:"-"`
	Code         string `orm:"column(code)" json:"code"`
	Name         string `orm:"column(name)" json:"name"`
	Note         string `orm:"column(note)" json:"note"`
	AuxData      int8   `orm:"column(aux_data)" json:"aux_data"`
	Abbreviation string `orm:"column(abbreviation)" json:"abbreviation"`
	Status       int8   `orm:"column(status)" json:"status"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *BusinessType) MarshalJSON() ([]byte, error) {
	type Alias BusinessType

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// // Save : function to save data into database. will update when it has valid id, create otherwise
// func (m *BusinessType) Save(fields ...string) (err error) {
// 	o := orm.NewOrm()
// 	if m.ID > 0 {
// 		_, err = o.Update(m, fields...)
// 	} else {
// 		m.ID, err = o.Insert(m)
// 	}
// 	return
// }

// Read : function to get data from database
func (m *BusinessType) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
