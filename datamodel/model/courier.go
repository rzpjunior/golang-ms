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
	orm.RegisterModel(new(Courier))
}

// Courier model for courier table.
type Courier struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	CourierID         int64     `orm:"-" json:"courier_id,omitempty"`
	Code              string    `orm:"column(code);size(50);null" json:"code"`
	Name              string    `orm:"column(name);size(100);null" json:"name"`
	PhoneNumber       string    `orm:"column(phone_number);size(50);null" json:"phone_number"`
	LicensePlate      string    `orm:"column(license_plate);size(15);null" json:"license_plate"`
	EmergencyMode     int8      `orm:"column(emergency_mode)" json:"emergency_mode"`
	LastEmergencyTime time.Time `orm:"column(last_emergency_time)" json:"last_emergency_time"`
	Status            int8      `orm:"column(status);null" json:"status"`

	Role            *Role            `orm:"column(role_id);null;rel(fk)" json:"role,omitempty"`
	User            *User            `orm:"column(user_id);null;rel(fk)" json:"user,omitempty"`
	VehicleProfiles *VehicleProfiles `orm:"column(vehicle_profiles_id);rel(fk);null" json:"vehicle_profiles,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Courier) MarshalJSON() ([]byte, error) {
	type Alias Courier

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating courier struct into courier table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to courier.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Courier) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting courier data
// this also will truncated all data from all table
// that have relation with this courier.
func (m *Courier) Delete() (err error) {
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
func (m *Courier) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
