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
	orm.RegisterModel(new(AuditLog))
}

// AuditLog model for audit_log table.
type AuditLog struct {
	ID           int64         `orm:"column(id);auto" json:"-"`
	Staff        *Staff        `orm:"column(staff_id);null;rel(fk)" json:"staff,omitempty"`
	UserMerchant *UserMerchant `orm:"column(merchant_id);null;rel(fk)" json:"user_merchant,omitempty"`
	RefID        int64         `orm:"column(ref_id);" json:"ref_id"`
	Type         string        `orm:"column(type);null" json:"type"`
	Function     string        `orm:"column(function);null" json:"function"`
	Timestamp    time.Time     `orm:"column(timestamp);type(timestamp);null" json:"timestamp"`
	Note         string        `orm:"column(note);size(250);null" json:"note"`

	ChangesLog ChangesLog `orm:"-" json:"changes_log,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *AuditLog) MarshalJSON() ([]byte, error) {
	type Alias AuditLog

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating AuditTrail struct into audit_trail table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to audit_trail.
// The field parameter is an field that will be saved, it is
// usefull for AuditLog updating data.
func (m *AuditLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting audit_trail data
// this also will truncated all data from all table
// that have relation with this audit_trail.
func (m *AuditLog) Delete() (err error) {
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
func (m *AuditLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
