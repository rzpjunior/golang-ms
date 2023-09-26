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
	orm.RegisterModel(new(Banner))
}

// Banner : struct to hold model data for database
type Banner struct {
	ID                 int64     `orm:"column(id);auto" json:"-"`
	Code               string    `orm:"column(code)" json:"code"`
	Name               string    `orm:"column(name)" json:"name"`
	Area               string    `orm:"column(area)" json:"area"`
	AreaName           string    `orm:"-" json:"area_name"`
	AreaNameArr        []string  `orm:"-" json:"area_name_arr"`
	Archetype          string    `orm:"column(archetype)" json:"archetype"`
	ArchetypeName      string    `orm:"-" json:"archetype_name"`
	ArchetypeNameArr   []string  `orm:"-" json:"archetype_name_arr"`
	StartDate          time.Time `orm:"column(start_date);type(timestamp);null" json:"start_date"`
	EndDate            time.Time `orm:"column(end_date);type(timestamp);null" json:"end_date"`
	NavigationType     int8      `orm:"column(navigate_type)" json:"navigate_type"`
	NavigationTypeName string    `orm:"-" json:"navigate_type_name"`
	NavigationUrl      string    `orm:"column(navigate_url);null" json:"navigate_url"`
	ImageUrl           string    `orm:"column(image_url)" json:"image_url"`
	Queue              int8      `orm:"column(queue)" json:"queue"`
	Status             int8      `orm:"column(status)" json:"status"`
	Note               string    `orm:"column(note)" json:"note"`

	TagProduct     *TagProduct     `orm:"column(tag_product_id);null;rel(fk)" json:"tag_product"`
	Product        *Product        `orm:"column(product_id);null;rel(fk)" json:"product"`
	ProductSection *ProductSection `orm:"column(product_section_id);null;rel(fk)" json:"product_section"`

	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy int64     `orm:"column(created_by)" json:"created_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Banner) MarshalJSON() ([]byte, error) {
	type Alias Banner

	return json.Marshal(&struct {
		ID        string `json:"id"`
		CreatedBy string `json:"created_by"`
		*Alias
	}{
		ID:        common.Encrypt(m.ID),
		CreatedBy: common.Encrypt(m.CreatedBy),
		Alias:     (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Banner) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Banner) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
