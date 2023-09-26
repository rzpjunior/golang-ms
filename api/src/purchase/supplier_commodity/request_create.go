// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_commodity

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	Code   string `json:"-"`
	Name   string `orm:"column(name)" json:"name" valid:"required"`
	Note   string `orm:"column(note)" json:"note"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("supplier_commodity"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if len(c.Name) > 100 {
		o.Failure("name.invalid", util.ErrorCharLength("name", 100))
	}

	exists, err := repository.IsExistsSupplierCommodity(c.Name)

	if err != nil {
		o.Failure("name", util.ErrorInvalidData("supplier commodity"))
	}

	if exists {
		o.Failure("name", util.ErrorDuplicate("supplier commodity"))
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":   util.ErrorInputRequired("name"),
	}
}
