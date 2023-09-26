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
	orm.RegisterModel(new(VoucherItem))
}

// VoucherItem model for voucher_item table.
type VoucherItem struct {
	ID           int64    `orm:"column(id);auto" json:"-"`
	Voucher      *Voucher `orm:"column(voucher_id);null;rel(fk)" json:"voucher,omitempty"`
	Product      *Product `orm:"column(product_id);null;rel(fk)"  json:"product,omitempty"`
	MinOrderDisc float64  `orm:"column(min_qty_disc);null;digits(5);decimals(2)" json:"min_qty_disc"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *VoucherItem) MarshalJSON() ([]byte, error) {
	type Alias VoucherItem

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into voucher_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *VoucherItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read execute select based on data struct that already
// assigned.
func (m *VoucherItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
