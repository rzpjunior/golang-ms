// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package prospect_supplier

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/api/datastore/repository"
)

type declineRequest struct {
	ID 		int64 `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`
}

func (c *declineRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if prospectsupplier, err := repository.ValidProspectSupplier(c.ID); err == nil {
		if prospectsupplier.RegStatus != 1 {
			o.Failure("reg_status.new", util.ErrorActive("reg_status"))
		}
	} else {
		o.Failure("prospect_supplier.invalid", util.ErrorInvalidData("prospect_supplier"))
	}

	return o
}

func (c *declineRequest) Messages() map[string]string {
	return map[string]string{}
}
