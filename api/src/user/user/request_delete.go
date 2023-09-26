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

type deleteRequest struct {
	ID           int64  `json:"-" valid:"required"`
	DeletionNote string `json:"note" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if user, err := repository.ValidUser(c.ID); err == nil {
		if user.Status != 2 {
			o.Failure("status.archived", util.ErrorArchived("status"))
		}

		staff := &model.Staff{User: user}
		staff.Read("User")

		if countStaff, err := repository.CountParentStaff(staff.ID); err == nil && countStaff > 0 {
			o.Failure("user.invalid", util.ErrorRelated("active and archived ", "staff", "user"))
		}
	} else {
		o.Failure("user.invalid", util.ErrorInvalidData("user"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *deleteRequest) Messages() map[string]string {
	return map[string]string{
		"deletion_note.required": util.ErrorInputRequired("deletion note"),
	}
}
