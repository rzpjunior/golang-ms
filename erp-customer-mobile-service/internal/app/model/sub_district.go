// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(AdmDivision))
}

// AdmDivision model for sub_district table.
type AdmDivision struct {
	ID             int64  `orm:"column(id);auto" json:"-"`
	Code           string `orm:"column(code);size(50);null" json:"code"`
	Value          string `orm:"column(value);size(50);null" json:"value"`
	Name           string `orm:"column(name);size(100);null" json:"name"`
	PostalCode     string `orm:"column(postal_code);size(10);null" json:"postal_code"`
	ConcatNoPrefix string `orm:"column(concat_no_prefix);size(250);null" json:"concat_no_prefix"`
	ConcatAddress  string `orm:"column(concat_address);size(250);null" json:"concat_address"`
	Note           string `orm:"column(note);size(250);null" json:"note"`
	Status         int8   `orm:"column(status);null" json:"status"`
	StatusConvert  string `orm:"-" json:"status_convert"`

	//District *District `orm:"column(district_id);null;rel(fk)" json:"district,omitempty"`
	Area *Region `orm:"-"  json:"region,omitempty"`
}
