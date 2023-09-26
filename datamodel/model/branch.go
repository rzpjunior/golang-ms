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
	orm.RegisterModel(new(Branch))
}

// Branch model for area table.
type Branch struct {
	ID                 int64     `orm:"column(id);auto" json:"-"`
	Code               string    `orm:"column(code)" json:"code,omitempty"`
	Name               string    `orm:"column(name)" json:"name,omitempty"`
	PicName            string    `orm:"column(pic_name)" json:"pic_name,omitempty"`
	PhoneNumber        string    `orm:"column(phone_number)" json:"phone_number,omitempty"`
	AltPhoneNumber     string    `orm:"column(alt_phone_number)" json:"alt_phone_number,omitempty"`
	AddressName        string    `orm:"column(address_name)" json:"address_name"`
	ShippingAddress    string    `orm:"column(shipping_address)" json:"shipping_address,omitempty"`
	Latitude           *float64  `orm:"column(latitude)" json:"latitude,omitempty"`
	Longitude          *float64  `orm:"column(longitude)" json:"longitude,omitempty"`
	PinpointValidation int8      `orm:"column(pinpoint_validation)" json:"pinpoint_validation"`
	Note               string    `orm:"column(note)" json:"note,omitempty"`
	MainBranch         int8      `orm:"column(main_branch)" json:"main_branch,omitempty"`
	Status             int8      `orm:"column(status)" json:"status"`
	CreatedAt          time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy          int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt      time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy      int64     `orm:"column(last_updated_by)" json:"last_updated_by"`

	Merchant    *Merchant    `orm:"column(merchant_id);null;rel(fk)" json:"merchant,omitempty"`
	Area        *Area        `orm:"column(area_id);null;rel(fk)" json:"area,omitempty"`
	Archetype   *Archetype   `orm:"column(archetype_id);null;rel(fk)" json:"archetype,omitempty"`
	PriceSet    *PriceSet    `orm:"column(price_set_id);null;rel(fk)" json:"price_set,omitempty"`
	Warehouse   *Warehouse   `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	Salesperson *Staff       `orm:"column(salesperson_id);null;rel(fk)" json:"salesperson,omitempty"`
	SubDistrict *SubDistrict `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`

	StatusConvert string `orm:"-" json:"status_convert"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Branch) MarshalJSON() ([]byte, error) {
	type Alias Branch

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into branch table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to branch.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Branch) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting branch data
// this also will truncated all data from all table
// that have relation with this user.
func (m *Branch) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *Branch) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
