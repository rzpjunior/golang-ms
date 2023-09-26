// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveHelperRequest struct {
	ID int64 `json:"-" valid:"required"`

	Staff   *model.Staff
	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *archiveHelperRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var count int
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	orSelect.Raw("select count(*) "+
		"from packing_order_item_assign poia inner join packing_order po inner join packing_order_item poi "+
		"where poia.packing_order_item_id = poi.id and poi.packing_order_id = po.id and poia.staff_id = ?", c.ID).QueryRow(&count)

	if count > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "helper", "packing order"))
	}
	if c.Staff, err = repository.ValidStaff(c.ID); err == nil {
		if c.Staff.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
		if c.Staff.User, err = repository.ValidUser(c.Staff.User.ID); err != nil {
			o.Failure("user.invalid", util.ErrorInvalidData("user"))
		}
	} else {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *archiveHelperRequest) Messages() map[string]string {
	return map[string]string{}
}
