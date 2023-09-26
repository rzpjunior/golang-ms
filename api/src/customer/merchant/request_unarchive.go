// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Merchant *model.Merchant `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 2 {
			o.Failure("status.archived", util.ErrorActive("status"))
		}
	} else {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
