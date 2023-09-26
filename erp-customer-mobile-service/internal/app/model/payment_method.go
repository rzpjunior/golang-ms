// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

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
