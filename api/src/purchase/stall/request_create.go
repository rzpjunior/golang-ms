// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stall

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold input create stall request data
type createRequest struct {
	Code        string `json:"-"`
	Name        string `json:"name" valid:"required|alpha_num_space|lte:100"`
	PhoneNumber string `json:"phone_number" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate create stall request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	c.Code = util.GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 3)

	// validate input length
	if len(c.PhoneNumber) < 8 || len(c.PhoneNumber) > 15 {
		o.Failure("phone_number", util.ErrorRangeChar("phone number", "8", "15"))
	}

	c.PhoneNumber = util.ParsePhoneNumberPrefix(c.PhoneNumber)

	// validate if phonenumber duplicate
	stall := &model.Stall{PhoneNumber: c.PhoneNumber}
	if err = stall.Read("PhoneNumber"); err == nil {
		o.Failure("phone_number", util.ErrorUnique("phone number"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":         util.ErrorInputRequired("name"),
		"phone_number.required": util.ErrorInputRequired("phone number"),
		"name.alpha_num_space":  util.ErrorAlphaNum("name"),
		"name.lte":              util.ErrorEqualLess("name", "100 characters"),
	}
}
