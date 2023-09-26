// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveHelperRequest struct {
	ID int64 `json:"-" valid:"required"`

	Staff   *model.Staff
	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *unarchiveHelperRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Staff, err = repository.ValidStaff(c.ID); err == nil {
		if c.Staff.Status != 2 {
			o.Failure("status.active", util.ErrorArchived("status"))
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
func (c *unarchiveHelperRequest) Messages() map[string]string {
	return map[string]string{}
}
