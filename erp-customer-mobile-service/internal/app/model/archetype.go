// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(Archetype))
}

// Archetype : struct to hold model data for database
type Archetype struct {
	ID               int64         `orm:"column(id);auto" json:"-"`
	CustomerGroup    int8          `orm:"column(customer_group)" json:"customer_group"`
	Code             string        `orm:"column(code)" json:"code"`
	Name             string        `orm:"column(name)" json:"name"`
	NameID           string        `orm:"column(name_id)" json:"name_id"`
	Abbreviation     string        `orm:"column(abbreviation)" json:"abbreviation"`
	Note             string        `orm:"column(note)" json:"note"`
	AuxData          int8          `orm:"column(aux_data)" json:"aux_data"`
	Status           int8          `orm:"column(status)" json:"status"`
	DocRequired      int64         `orm:"-" json:"document_required"`
	DocImageRequired []string      `orm:"-" json:"document_image_required"`
	BusinessType     *CustomerType `orm:"-" json:"customer_type"`
}
