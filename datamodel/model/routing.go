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
	orm.RegisterModel(new(Routing))
}

// Routing: struct to hold model data for database
type Routing struct {
	ID                          int64     `orm:"column(id);auto" json:"-"`
	Code                        string    `orm:"column(code)" json:"code"`
	RoutingGoal                 string    `orm:"column(routing_goal)" json:"routing_goal"`
	Wrt                         string    `orm:"column(wrt)" json:"wrt"`
	AvailableVehicles           string    `orm:"column(available_vehicles)" json:"-"`
	TotalSalesOrder             int64     `orm:"column(total_sales_order)" json:"total_sales_order"`
	RoutedSalesOrder            int64     `orm:"column(routed_sales_order)" json:"routed_sales_order"`
	DroppedSalesOrder           int64     `orm:"column(dropped_sales_order)" json:"dropped_sales_order"`
	TotalWeight                 float64   `orm:"column(total_weight)" json:"total_weight"`
	TotalFragileWeight          float64   `orm:"column(total_fragile_weight)" json:"total_fragile_weight"`
	TotalKoli                   float64   `orm:"column(total_koli)" json:"total_koli"`
	TotalCost                   float64   `orm:"column(total_cost)" json:"total_cost"`
	TotalBranch                 int64     `orm:"column(total_branch)" json:"total_branch"`
	DeliveryDate                time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	Status                      int8      `orm:"column(status)" json:"status"`
	ErrorResponse               string    `orm:"column(error_response)" json:"error_response"`
	MultiBatch                  int8      `orm:"column(multi_batch)" json:"multi_batch"`
	Priority                    int8      `orm:"column(priority)" json:"priority"`
	IncludeCourierAssignedOrder int8      `orm:"column(include_courier_assigned_order)" json:"include_courier_assigned_order"`
	KeepCourierAssignment       int8      `orm:"column(keep_courier_assignment)" json:"keep_courier_assignment"`
	ServiceTime                 int64     `orm:"column(service_time)" json:"service_time"`
	SetupTime                   int64     `orm:"column(setup_time)" json:"setup_time"`
	EnableCod                   int8      `orm:"column(enable_cod)" json:"enable_cod"`
	CreatedAt                   time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	ResponseAt                  time.Time `orm:"column(response_at)" json:"response_at"`
	CreatedBy                   *Staff    `orm:"column(created_by);rel(fk)" json:"created_by"`

	Warehouse               *Warehouse                 `orm:"column(warehouse_id);rel(fk)" json:"warehouse"`
	VehicleAvailablesDetail []*VehicleAvailablesDetail `orm:"-" json:"available_vehicles"`
}

type VehicleQuantityDetails struct {
	ID          int64   `json:"id"`
	Qty         int64   `json:"qty"`
	Routed      int64   `json:"routed"`
	SpeedFactor float64 `json:"speed_factor"`
}

type VehicleAvailablesDetail struct {
	ID             int64   `json:"id"`
	Qty            int64   `json:"qty"`
	Routed         int64   `json:"routed"`
	Name           string  `json:"name"`
	RoutingProfile string  `json:"routing_profile"`
	SpeedFactor    float64 `json:"speed_factor"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Routing) MarshalJSON() ([]byte, error) {
	type Alias Routing

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating routing struct into routing table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to routing.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Routing) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting routing data
// this also will truncated all data from all table
// that have relation with this routing .
func (m *Routing) Delete() (err error) {
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
func (m *Routing) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
