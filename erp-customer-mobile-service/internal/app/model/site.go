// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(Site))
}

// Site model for Site table.
type Site struct {
	ID             int64  `orm:"column(id);auto" json:"-"`
	Code           string `orm:"column(code);size(50);null" json:"code"`
	Name           string `orm:"column(name);size(100);null" json:"name"`
	PicName        string `orm:"column(pic_name);size(100);null" json:"pic_name"`
	PhoneNumber    string `orm:"column(phone_number);size(15);null" json:"phone_number"`
	AltPhoneNumber string `orm:"column(alt_phone_number);size(100);null" json:"alt_phone_number"`
	StreetAddress  string `orm:"column(street_address);size(350);null" json:"street_address"`
	Note           string `orm:"column(note);size(250);null" json:"note"`
	AuxData        int8   `orm:"column(aux_data)" json:"aux_data"`
	//MainSite  int8   `orm:"column(main_Site);null" json:"main_Site"`
	Status int8 `orm:"column(status);null" json:"status"`

	Region      *Region      `orm:"-" json:"region,omitempty"`
	AdmDivision *AdmDivision `orm:"-" json:"adm_division,omitempty"`

	StatusConvert string `orm:"-" json:"status_convert"`
}
