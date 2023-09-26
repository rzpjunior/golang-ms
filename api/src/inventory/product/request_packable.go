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

type packableRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *packableRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if product, err := repository.ValidProduct(c.ID); err == nil {
		if product.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}

		if product.Packability != 2 {
			o.Failure("packability.invalid", util.ErrorActiveIsPackable("not packable"))
		}
	} else {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
	}

	return o
}

func (c *packableRequest) Messages() map[string]string {
	return map[string]string{}
}
