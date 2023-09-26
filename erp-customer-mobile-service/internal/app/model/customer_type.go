// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(CustomerType))
}

// CustomerType : struct to hold model data for database
type CustomerType struct {
	ID               int64  `orm:"column(id);auto" json:"-"`
	Code             string `orm:"column(code)" json:"code"`
	Name             string `orm:"column(name)" json:"name"`
	Note             string `orm:"column(note)" json:"note"`
	AuxData          int8   `orm:"column(aux_data)" json:"aux_data"`
	Status           int8   `orm:"column(status)" json:"status"`
	DocImageRequired string `orm:"column(doc_image_required)" json:"doc_image_required"`
}
