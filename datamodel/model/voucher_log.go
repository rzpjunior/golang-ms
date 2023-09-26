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
	orm.RegisterModel(new(VoucherLog))
}

// VoucherLog model for voucher table.
type VoucherLog struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	TagCustomer       string    `orm:"column(tag_customer)" json:"customer_tag"`
	VoucherDiscAmount float64   `orm:"column(vou_disc_amount)" json:"voucher_discount_amount"`
	Timestamp         time.Time `orm:"column(timestamp)" json:"timestamp"`
	Status            int8      `orm:"column(status)" json:"status"`

	Voucher    *Voucher    `orm:"column(voucher_id);null;rel(fk)" json:"voucher"`
	Merchant   *Merchant   `orm:"column(merchant_id);null;rel(fk)" json:"merchant"`
	Branch     *Branch     `orm:"column(branch_id);null;rel(fk)" json:"branch"`
	SalesOrder *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *VoucherLog) MarshalJSON() ([]byte, error) {
	type Alias VoucherLog

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

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *VoucherLog) Save(fields ...string) (err error) {
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
func (m *VoucherLog) Delete() (err error) {
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
func (m *VoucherLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
