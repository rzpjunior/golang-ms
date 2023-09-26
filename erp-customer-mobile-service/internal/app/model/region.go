// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(Region))
}

// Region model for Region table.
type Region struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	Code          string `orm:"column(code);size(50);null" json:"code"`
	Value         string `orm:"column(value);size(250);null" json:"value"`
	Name          string `orm:"column(name);size(100);null" json:"name"`
	Note          string `orm:"column(note);size(250);null" json:"note"`
	Status        int8   `orm:"column(status);null" json:"status"`
	AuxData       int8   `orm:"column(aux_data)" json:"aux_data"`
	StatusConvert string `orm:"-" json:"status_convert"`
}
