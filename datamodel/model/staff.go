// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(Staff))
}

// Staff model for staff table.
type Staff struct {
	ID                 int64        `orm:"column(id);auto" json:"-"`
	StaffID            int64        `orm:"-" json:"staff_id,omitempty"`
	Code               string       `orm:"column(code);size(50);null" json:"code"`
	Name               string       `orm:"column(name);size(100);null" json:"name"`
	DisplayName        string       `orm:"column(display_name);size(100);null" json:"display_name"`
	EmployeeCode       string       `orm:"column(employee_code);size(50);null" json:"employee_code"`
	RoleGroup          int8         `orm:"column(role_group);null" json:"role_group"`
	PhoneNumber        string       `orm:"column(phone_number);size(15);null" json:"phone_number"`
	Status             int8         `orm:"column(status);null" json:"status"`
	SalesGroupID       int64        `orm:"column(sales_group_id);null;" json:"sales_group_id,omitempty"`
	SalesGroupName     string       `orm:"-" json:"sales_group_name,omitempty"`
	WarehouseAccessStr string       `orm:"column(warehouse_access)" json:"warehouse_access_str"`
	StatusConvert      string       `orm:"-" json:"status_convert"`
	WarehouseAccess    []*Warehouse `orm:"-" json:"warehouse_access"`

	// Picking List Module
	UsedStaff bool `orm:"-" json:"used_staff"`
	IsBusy    bool `orm:"-" json:"is_busy"`

	Role      *Role      `orm:"column(role_id);null;rel(fk)" json:"role,omitempty"`
	User      *User      `orm:"column(user_id);null;rel(fk)" json:"user,omitempty"`
	Area      *Area      `orm:"column(area_id);null;rel(fk)" json:"area,omitempty"`
	Parent    *Staff     `orm:"column(parent_id);null;rel(fk)" json:"parent,omitempty"`
	Warehouse *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted
// .
func (m *Staff) MarshalJSON() ([]byte, error) {
	type Alias Staff

	alias := &struct {
		ID            string `json:"id"`
		SalesGroupID  string `json:"sales_group_id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		SalesGroupID:  common.Encrypt(m.SalesGroupID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Staff) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *Staff) Delete() (err error) {
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
func (m *Staff) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
