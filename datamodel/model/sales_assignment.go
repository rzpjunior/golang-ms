// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(SalesAssignment))
}

// Sales Assignment: struct to hold model data for database
type SalesAssignment struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	Code      string    `orm:"column(code);size(50);null" json:"code"`
	StartDate time.Time `orm:"column(start_date)" json:"start_date"`
	EndDate   time.Time `orm:"column(end_date)" json:"end_date"`
	Status    int8      `orm:"column(status)" json:"status"`

	SalesGroup *SalesGroup `orm:"column(sales_group_id);null;rel(fk)" json:"sales_group"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesAssignment) MarshalJSON() ([]byte, error) {
	type Alias SalesAssignment

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusDoc(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesAssignment) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesAssignment) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
