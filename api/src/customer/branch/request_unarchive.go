// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Branch   *model.Branch
	Merchant *model.Merchant

	Session *auth.SessionData
}

// Validate : function to validate request data
func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Branch, err = repository.ValidBranch(c.ID); err == nil {
		if c.Branch.Status != 2 {
			o.Failure("status.archive", util.ErrorArchived("status"))
		} else {
			c.Merchant = &model.Merchant{ID: c.Branch.Merchant.ID}
			c.Merchant.Read("ID")
		}
	} else {
		o.Failure("branch.invalid", util.ErrorInvalidData("branch"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
