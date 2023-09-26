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
	orm.RegisterModel(new(CreditLimitLog))
}

// CreditLimitLog model for audit_log table.
type CreditLimitLog struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	Merchant          int64     `orm:"column(merchant_id);null" json:"merchant,omitempty"`
	RefID             int64     `orm:"column(ref_id);" json:"ref_id"`
	Type              string    `orm:"column(type);null" json:"type"`
	CreditLimitBefore float64   `orm:"column(credit_limit_before)" json:"credit_limit_before"`
	CreditLimitAfter  float64   `orm:"column(credit_limit_after)" json:"credit_limit_after"`
	Note              string    `orm:"column(note);size(250);null" json:"note"`
	CreatedAt         time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy         int64     `orm:"column(created_by)" json:"created_by"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *CreditLimitLog) MarshalJSON() ([]byte, error) {
	type Alias CreditLimitLog

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
// usefull for CreditLimitLog updating data.
func (m *CreditLimitLog) Save(fields ...string) (err error) {
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
func (m *CreditLimitLog) Delete() (err error) {
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
func (m *CreditLimitLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
