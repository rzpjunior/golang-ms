// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type notFragileRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *notFragileRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var product *model.Product
	var err error

	if product, err = repository.ValidProduct(c.ID); err != nil {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
	}

	if product.FragileGoods != 1 {
		o.Failure("fragility.invalid", util.ErrorActiveIsPackable("fragile"))
	}

	return o
}

func (c *notFragileRequest) Messages() map[string]string {
	return map[string]string{}
}
