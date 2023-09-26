// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package method

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold payment method request data
type createRequest struct {
	Code string `json:"-"`
	Name string `json:"name" valid:"required"`
	Note string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate payment method request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("payment_method"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	filter := map[string]interface{}{"name": c.Name}
	exclude := map[string]interface{}{}
	if _, countName, err := repository.CheckPaymentMethodData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required": util.ErrorInputRequired("name"),
	}
}
