// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if product, err := repository.ValidProduct(c.ID); err == nil {
		if product.Status != 2 {
			o.Failure("status.archive", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
	}

	return o
}

func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
