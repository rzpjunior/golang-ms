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
	orm.RegisterModel(new(PaymentGroup))
}

// PaymentGroup : struct to hold payment group sls model data for database
type PaymentGroup struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	NameID string `orm:"column(name_id)" json:"name_id"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status)" json:"status"`
}

// TableName : set table name used by model
func (PaymentGroup) TableName() string {
	return "payment_group_sls"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PaymentGroup) MarshalJSON() ([]byte, error) {
	type Alias PaymentGroup

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *PaymentGroup) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
