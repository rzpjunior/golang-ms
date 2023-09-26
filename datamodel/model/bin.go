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
	orm.RegisterModel(new(Bin))
}

// Bin model for bin table.
type Bin struct {
	ID          int64      `orm:"column(id);auto" json:"-"`
	Code        string     `orm:"column(code)" json:"code"`
	Name        string     `orm:"column(name)" json:"name"`
	Warehouse   *Warehouse `orm:"column(warehouse_id);rel(fk)" json:"warehouse"`
	Product     *Product   `orm:"column(product_id);rel(fk)" json:"product"`
	Latitude    *float64   `orm:"column(latitude)" json:"latitude"`
	Longitude   *float64   `orm:"column(longitude)" json:"longitude"`
	Status      int8       `orm:"column(status)" json:"status"`
	Note        string     `orm:"column(note)" json:"note"`
	ServiceTime int64      `orm:"column(service_time)" json:"service_time"`
	CreatedAt   time.Time  `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy   *Staff     `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	UpdatedAt   time.Time  `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	UpdatedBy   *Staff     `orm:"column(updated_by);null;rel(fk)" json:"updated_by,omitempty"`

	ProductName string `orm:"-" json:"product_name"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Bin) MarshalJSON() ([]byte, error) {
	type Alias Bin

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Bin) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Bin) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
