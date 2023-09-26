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
	orm.RegisterModel(new(PaymentChannel))
}

// PaymentChannel : struct to hold payment method model data for database
type PaymentChannel struct {
	ID         int64  `orm:"column(id);auto" json:"-"`
	Code       string `orm:"column(code);size(50);null" json:"code"`
	Value      string `orm:"column(value);size(50);null" json:"value"`
	Name       string `orm:"column(name);size(100);null" json:"name"`
	ImageUrl   string `orm:"column(image_url);size(300);null" json:"image_url"`
	Note       string `orm:"column(note);size(255)" json:"note"`
	Status     int8   `orm:"column(status)" json:"status"`
	PublishIva int8   `orm:"column(publish_iva)" json:"publish_iva"`
	PublishFva int8   `orm:"column(publish_fva)" json:"publish_fva"`

	PaymentMethod *PaymentMethod `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PaymentChannel) MarshalJSON() ([]byte, error) {
	type Alias PaymentChannel

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *PaymentChannel) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
