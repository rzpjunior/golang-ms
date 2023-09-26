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
	orm.RegisterModel(new(VehicleProfiles))
}

// VehicleProfiles : struct to hold model data for database
type VehicleProfiles struct {
	ID             int64   `orm:"column(id);auto" json:"-"`
	Code           string  `orm:"column(code)" json:"code"`
	Name           string  `orm:"column(name)" json:"name"`
	MaxKoli        float32 `orm:"column(max_koli)" json:"max_koli"`
	MaxWeight      float32 `orm:"column(max_weight)" json:"max_weight"`
	MaxFragile     float32 `orm:"column(max_fragile)" json:"max_fragile"`
	SpeedFactor    float32 `orm:"column(speed_factor)" json:"speed_factor"`
	Skills         string  `orm:"column(skills)" json:"skills"`
	Status         int8    `orm:"column(status)" json:"status"`
	InitialCost    float64 `orm:"column(initial_cost)" json:"initial_cost"`
	SubsequentCost float64 `orm:"column(subsequent_cost)" json:"subsequent_cost"`
	MaxVehicle     int64   `orm:"column(max_available_vehicle)" json:"max_available_vehicle"`

	CourierVendor  *CourierVendor `orm:"column(courier_vendor_id);rel(fk)" json:"courier_vendor"`
	RoutingProfile *Glossary      `orm:"column(routing_profile);rel(fk)" json:"routing_profile"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *VehicleProfiles) MarshalJSON() ([]byte, error) {
	type Alias VehicleProfiles

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *VehicleProfiles) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *VehicleProfiles) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
