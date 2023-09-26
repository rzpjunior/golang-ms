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
	orm.RegisterModel(new(SalesGroup))
}

// Sales Group: struct to hold model data for database
type SalesGroup struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	Name             string    `orm:"column(name);null" json:"name"`
	City             string    `orm:"column(city)" json:"-"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy        int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt    time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy    int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
	Status           int8      `orm:"column(status)" json:"status"`
	CityStr          string    `orm:"-" json:"city_str"`
	SalespersonTotal int64     `orm:"-" json:"salesperson_total"`

	BusinessType    *BusinessType     `orm:"column(business_type_id);null;rel(fk)" json:"business_type"`
	SalesManager    *Staff            `orm:"column(sls_man_id);null;rel(fk)" json:"sls_man"`
	Area            *Area             `orm:"column(area_id);null;rel(fk)" json:"area"`
	SalesGroupItems []*SalesGroupItem `orm:"reverse(many)" json:"sales_group_item,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesGroup) MarshalJSON() ([]byte, error) {
	type Alias SalesGroup

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesGroup) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesGroup) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
