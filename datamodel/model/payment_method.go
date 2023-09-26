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
	orm.RegisterModel(new(PaymentMethod))
}

// PaymentMethod : struct to hold payment method model data for database
type PaymentMethod struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Code        string `orm:"column(code)" json:"code"`
	Name        string `orm:"column(name)" json:"name"`
	Note        string `orm:"column(note)" json:"note"`
	Status      int8   `orm:"column(status)" json:"status"`
	Publish     int8   `orm:"column(publish)" json:"publish"`
	Maintenance int8   `orm:"column(maintenance)" json:"maintenance"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PaymentMethod) MarshalJSON() ([]byte, error) {
	type Alias PaymentMethod

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *PaymentMethod) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
