// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_commodity

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID     int64  `json:"-"`
	Name   string `orm:"column(name)" json:"name" valid:"required"`
	Note   string `orm:"column(note)" json:"note"`

	SupplierCommodity *model.SupplierCommodity `json:"-"`
	Session           *auth.SessionData        `json:"-"`
}

// Validate : function to validate request data
func (u *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if len(u.Name) > 100 {
		o.Failure("name.invalid", util.ErrorCharLength("name", 100))
	}

	u.SupplierCommodity, err = repository.GetSupplierCommodity("id", u.ID)

	if err != nil {
		o.Failure("id", util.ErrorInvalidData("supplier commodity"))
	}

	if u.Name != u.SupplierCommodity.Name {

		exists, err := repository.IsExistsSupplierCommodity(u.Name)

		if err != nil {
			o.Failure("name", util.ErrorInvalidData("supplier commodity"))
		}
	
		if exists {
			o.Failure("name", util.ErrorDuplicate("supplier commodity"))
		}

	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required": util.ErrorInputRequired("name"),
	}
}
