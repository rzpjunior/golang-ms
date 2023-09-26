// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_organization

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateRequest struct {
	ID      int64  `json:"-" valid:"required"`
	Name    string `json:"name" valid:"required|alpha_num_space|lte:100"`
	Address string `json:"address" valid:"required|lte:350"`
	Note    string `json:"note" valid:"lte:250"`

	SupplierOrganization *model.SupplierOrganization `json:"-"`
	Session              *auth.SessionData           `json:"-"`
}

// Validate : function to validate request data
func (u *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	u.SupplierOrganization, err = repository.ValidSupplierOrganization(u.ID)
	if err != nil {
		o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization id"))
	}

	if u.Name != u.SupplierOrganization.Name {

		exists, err := repository.IsExistsSupplierOrganization(u.Name)

		if err != nil {
			o.Failure("name", util.ErrorInvalidData("supplier organization"))
		}

		if exists {
			o.Failure("name", util.ErrorDuplicate("supplier organization"))
		}

	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"id.required":          util.ErrorInputRequired("id"),
		"name.required":        util.ErrorInputRequired("name"),
		"name.alpha_num_space": util.ErrorAlphaNum("name"),
		"name.lte":             util.ErrorEqualLess("name", "100"),
		"address.required":     util.ErrorInputRequired("address"),
		"address.lte":          util.ErrorEqualLess("address", "350"),
		"note.lte":             util.ErrorEqualLess("note", "250"),
	}
}
