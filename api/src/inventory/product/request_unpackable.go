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

type unpackableRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *unpackableRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var total int64
	var filter, exclude map[string]interface{}

	if product, err := repository.ValidProduct(c.ID); err == nil {
		if product.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}

		if product.Packability != 1 {
			o.Failure("packability.invalid", util.ErrorActiveIsPackable("packable"))
		}

		filter = map[string]interface{}{"PackingOrder__status": 1, "Product__id": c.ID}
		if _, total, err = repository.CheckPackingOrderItemData(filter, exclude); err == nil && total > 0 {
			o.Failure("packing_o", util.ErrorExistActivePackingOrder())
		}
	} else {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
	}

	return o
}

func (c *unpackableRequest) Messages() map[string]string {
	return map[string]string{}
}
