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
	orm.RegisterModel(new(UserMerchant))
}

// UserMerchant model for user_merchant table.
type UserMerchant struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	Code          string `orm:"column(code);size(50);null" json:"code,omitempty"`
	LoginToken    string `orm:"column(login_token);size(100);null" json:"login_token,omitempty"`
	FirebaseToken string `orm:"column(firebase_token);size(250);null" json:"firebase_token,omitempty"`
	//FirebaseID    string    `orm:"column(firebase_id);size(250);null" json:"firebase_id,omitempty"`
	Verification  int8      `orm:"column(verification)" json:"verification,omitempty"`
	TncAccVersion string    `orm:"column(tnc_acc_version);size(50);null" json:"tnc_acc_version,omitempty"`
	TncAccAt      time.Time `orm:"column(tnc_acc_at);type(timestamp);null" json:"tnc_acc_at"`
	LastLoginAt   time.Time `orm:"column(last_login_at);type(timestamp);null" json:"last_login_at"`
	Note          string    `orm:"column(note);size(250);null" json:"note,omitempty"`
	Status        int8      `orm:"column(status);null" json:"status,omitempty"`
	ForceLogout   int8      `orm:"column(force_logout);null" json:"force_logout,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *UserMerchant) MarshalJSON() ([]byte, error) {
	type Alias UserMerchant

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
	//if m.Password != "" {
	//	m.Password = "********"
	//}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *UserMerchant) Save(fields ...string) (err error) {
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
func (m *UserMerchant) Delete() (err error) {
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
func (m *UserMerchant) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
