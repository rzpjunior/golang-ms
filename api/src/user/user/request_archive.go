// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if user, err := repository.ValidUser(c.ID); err == nil {
		if user.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("user.invalid", util.ErrorInvalidData("user"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
