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
	orm.RegisterModel(new(User))
}

// User model for user table.
type User struct {
	ID                  int64     `orm:"column(id);auto" json:"-"`
	Code                string    `orm:"column(code);size(50);null" json:"code"`
	Email               string    `orm:"column(email);size(100);null" json:"email"`
	Password            string    `orm:"column(password);size(250);null" json:"password"`
	Status              int8      `orm:"column(status);null" json:"status"`
	PickingNotifToken   string    `orm:"column(picking_notif_token);null" json:"picking_notif_token"`
	LastLoginAt         time.Time `orm:"column(last_login_at);type(timestamp);null" json:"last_login_at"`
	Note                string    `orm:"column(note);size(250);null" json:"note"`
	ForceLogout         int8      `orm:"column(force_logout);null" json:"force_logout"`
	DashboardNotifToken string    `orm:"column(dashboard_notif_token);null" json:"dashboard_notif_token,omitempty"`
	SalesAppLoginToken  string    `orm:"column(salesapp_login_token);null" json:"salesapp_login_token,omitempty"`
	SalesAppNotifToken  string    `orm:"column(salesapp_notif_token);size(250);null" json:"salesapp_notif_token,omitempty"`
	PurchaserNotifToken string    `orm:"column(purchaser_notif_token);size(250);null" json:"purchaser_notif_token,omitempty"`
	StatusConvert       string    `orm:"-" json:"status_convert"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *User) MarshalJSON() ([]byte, error) {
	type Alias User

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	// hide password user
	if m.Password != "" {
		m.Password = "********"
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *User) Save(fields ...string) (err error) {
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
func (m *User) Delete() (err error) {
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
func (m *User) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
