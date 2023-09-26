// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(Staff))
}

// Staff model for staff table.
type Staff struct {
	ID            int64   `orm:"column(id);auto" json:"-"`
	Role          *Role   `orm:"-" json:"role,omitempty"`
	User          *User   `orm:"-"  json:"user,omitempty"`
	Region        *Region `orm:"-" json:"area,omitempty"`
	Parent        *Staff  `orm:"-" json:"parent,omitempty"`
	Site          *Site   `orm:"-" json:"site,omitempty"`
	Code          string  `orm:"column(code);size(50);null" json:"code"`
	Name          string  `orm:"column(name);size(100);null" json:"name"`
	DisplayName   string  `orm:"column(display_name);size(100);null" json:"display_name"`
	EmployeeCode  string  `orm:"column(employee_code);size(50);null" json:"employee_code"`
	RoleGroup     int8    `orm:"column(role_group);null" json:"role_group"`
	PhoneNumber   string  `orm:"column(phone_number);size(15);null" json:"phone_number"`
	Status        int8    `orm:"column(status);null" json:"status"`
	SalesGroupID  int64   `orm:"column(sales_group_id);null;" json:"sales_group_id,omitempty"`
	StatusConvert string  `orm:"-" json:"status_convert"`
}
