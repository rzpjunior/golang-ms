// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_type

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// update : struct to hold supplier type request data
type updateRequest struct {
	ID                    int64    `json:"-"`
	Code                  string   `json:"-"`
	Name                  string   `json:"name" valid:"required"`
	Abbreviation          string   `json:"abbreviation" valid:"required"`
	Note                  string   `json:"note"`

	SupplierType        *model.SupplierType        `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier type request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if len(c.Name) > 100 {
		o.Failure("name.invalid", util.ErrorCharLength("name", 100))
	}

	if len(c.Abbreviation) > 3 {
		o.Failure("abbreviation.invalid", util.ErrorCharLength("abbreviation", 3))
	}

	c.SupplierType, err = repository.GetSupplierType("id", c.ID)

	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))
		return o
	}

	if c.Name != c.SupplierType.Name {
		exists, err := repository.IsExistsSupplierType(c.Name)

		if err != nil {
			o.Failure("name", util.ErrorInvalidData("supplier type"))
		}
	
		if exists {
			o.Failure("name", util.ErrorDuplicate("supplier type"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":                   util.ErrorInputRequired("name"),
		"abbreviation.required":           util.ErrorInputRequired("abbreviation"),
	}
}
