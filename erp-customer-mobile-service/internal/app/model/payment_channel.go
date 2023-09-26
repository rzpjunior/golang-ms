// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(PaymentChannel))
}

// PaymentChannel : struct to hold payment method model data for database
type PaymentChannel struct {
	ID              int64  `orm:"column(id);auto" json:"-"`
	Code            string `orm:"column(code);size(50);null" json:"code"`
	Value           string `orm:"column(value);size(50);null" json:"value"`
	Name            string `orm:"column(name);size(100);null" json:"name"`
	ImageUrl        string `orm:"column(image_url);size(300);null" json:"image_url"`
	PaymentGuideUrl string `orm:"column(payment_guide_url);size(300);null" json:"payment_guide_url"`
	Note            string `orm:"column(note);size(255)" json:"note"`
	Status          int8   `orm:"column(status)" json:"status"`
	PublishIva      int8   `orm:"column(publish_iva)" json:"publish_iva"`
	PublishFva      int8   `orm:"column(publish_fva)" json:"publish_fva"`

	PaymentMethod *PaymentMethod `orm:"-"  json:"payment_method"`
}
