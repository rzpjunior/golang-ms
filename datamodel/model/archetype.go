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
	orm.RegisterModel(new(Archetype))
}

// Archetype : struct to hold model data for database
type Archetype struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	CustomerGroup int8   `orm:"column(customer_group)" json:"customer_group"`
	Code          string `orm:"column(code)" json:"code"`
	Name          string `orm:"column(name)" json:"name"`
	NameID        string `orm:"column(name_id)" json:"name_id"`
	Abbreviation  string `orm:"column(abbreviation)" json:"abbreviation"`
	Note          string `orm:"column(note)" json:"note"`
	AuxData       int8   `orm:"column(aux_data)" json:"aux_data"`
	Status        int8   `orm:"column(status)" json:"status"`

	BusinessType *BusinessType `orm:"column(business_type_id);null;rel(fk)" json:"business_type"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Archetype) MarshalJSON() ([]byte, error) {
	type Alias Archetype

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Archetype) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Archetype) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
