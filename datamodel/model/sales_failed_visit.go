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
	orm.RegisterModel(new(SalesFailedVisit))
}

// SalesFailedVisit: struct to hold model data for database
type SalesFailedVisit struct {
	ID                  int64                `orm:"column(id);auto" json:"-"`
	SalesAssignmentItem *SalesAssignmentItem `orm:"column(sales_assignment_item_id);null;rel(fk)" json:"sales_assignment_item"`
	FailedStatus        int8                 `orm:"column(failed_status)" json:"failed_status"`
	DescriptionFailed   string               `orm:"column(description_failed)" json:"description_failed"`
	FailedImage         string               `orm:"column(failed_image)" json:"failed_image"`

	FailedImageList []string `orm:"-" json:"failed_image_list"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesFailedVisit) MarshalJSON() ([]byte, error) {
	type Alias SalesFailedVisit

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Read : function to get data from database
func (m *SalesFailedVisit) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
