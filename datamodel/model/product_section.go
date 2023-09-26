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
	orm.RegisterModel(new(ProductSection))
}

// ProductSection : struct to hold model data for database
type ProductSection struct {
	ID                 int64                 `orm:"column(id);auto" json:"-"`
	Code               string                `orm:"column(code)" json:"code"`
	Name               string                `orm:"column(name)" json:"name"`
	BackgroundImage    string                `orm:"column(background_image)" json:"background_image"`
	Area               string                `orm:"column(area)" json:"area"`
	AreaArr            []*Area               `orm:"-" json:"area_arr,omitempty"`
	AreaName           string                `orm:"-" json:"area_name"`
	AreaNameArr        []string              `orm:"-" json:"area_name_arr"`
	Archetype          string                `orm:"column(archetype)" json:"archetype"`
	ArchetypeArr       []*Archetype          `orm:"-" json:"archetype_arr,omitempty"`
	ArchetypeName      string                `orm:"-" json:"archetype_name"`
	ArchetypeNameArr   []string              `orm:"-" json:"archetype_name_arr"`
	StartAt            time.Time             `orm:"column(start_at);type(timestamp);null" json:"start_at"`
	EndAt              time.Time             `orm:"column(end_at);type(timestamp);null" json:"end_at"`
	Sequence           int8                  `orm:"column(sequence)" json:"sequence"`
	Status             int8                  `orm:"column(status)" json:"status"`
	Note               string                `orm:"column(note)" json:"note"`
	Product            string                `orm:"column(product)" json:"product"`
	CreatedAt          time.Time             `orm:"column(created_at)" json:"created_at"`
	UpdatedAt          time.Time             `orm:"column(updated_at)" json:"updated_at"`
	ProductSectionItem []*ProductSectionItem `orm:"-" json:"product_section_item,omitempty"`
	Type               int8                  `orm:"column(type)" json:"type"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ProductSection) MarshalJSON() ([]byte, error) {
	type Alias ProductSection

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ProductSection) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ProductSection) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
