// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

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
