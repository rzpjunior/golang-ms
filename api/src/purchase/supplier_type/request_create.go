// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_type

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold supplier type request data
type createRequest struct {
	Code                  string   `json:"-"`
	Name                  string   `json:"name" valid:"required"`
	Abbreviation          string   `json:"abbreviation" valid:"required"`
	Note                  string   `json:"note"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier type request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("supplier_type"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if len(c.Name) > 100 {
		o.Failure("name.invalid", util.ErrorCharLength("name", 100))
	}

	if len(c.Abbreviation) > 3 {
		o.Failure("abbreviation.invalid", util.ErrorCharLength("abbreviation", 3))
	}


	exists, err := repository.IsExistsSupplierType(c.Name)

	if err != nil {
		o.Failure("name", util.ErrorInvalidData("supplier type"))
	}

	if exists {
		o.Failure("name", util.ErrorDuplicate("supplier type"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":                   util.ErrorInputRequired("name"),
		"abbreviation.required":           util.ErrorInputRequired("abbreviation"),
	}
}
