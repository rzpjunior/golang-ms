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
	orm.RegisterModel(new(AreaBusinessPolicy))
}

// AreaBusinessPolicy : struct to hold model data for database
type AreaBusinessPolicy struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	MinOrder      float64   `orm:"column(min_order);null;digits(20);decimals(2)" json:"min_order"`
	DeliveryFee   float64   `orm:"column(delivery_fee);null;digits(20);decimals(2)" json:"delivery_fee"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy int64     `orm:"column(last_updated_by)" json:"last_updated_by"`

	Area         *Area         `orm:"column(area_id);rel(fk)" json:"area,omitempty"`
	BusinessType *BusinessType `orm:"column(business_type_id);rel(fk)" json:"business,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *AreaBusinessPolicy) MarshalJSON() ([]byte, error) {
	type Alias AreaBusinessPolicy

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *AreaBusinessPolicy) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *AreaBusinessPolicy) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
