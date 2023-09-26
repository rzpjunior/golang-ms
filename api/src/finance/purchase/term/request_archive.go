// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package term

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

	if purchaseTerm, err := repository.ValidPurchaseTerm(c.ID); err == nil {
		if purchaseTerm.Status != 1 {
			o.Failure("id.invalid", util.ErrorActive("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("purchase term"))
	}

	if countSupplier, e := repository.CountNonDeletedSupplierByPurchaseTermId(c.ID); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("purchase term"))
	} else {
		if countSupplier > int64(0) {
			o.Failure("id.invalid", util.ErrorRelated("active and archived ", "supplier", "payment term"))
		}
	}

	return o
}

// Messages : function to return error messages after validation
func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
