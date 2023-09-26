// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type suspensionRequest struct {
	MerchantCode       string `json:"merchant_code" valid:"required"`
	Merchant           *model.Merchant
	Session            *auth.SessionData `json:"-"`
	SuspendOrUnSuspend string            `json:"-"`
}

// Validate : function to validate request data
func (c *suspensionRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	c.Merchant = &model.Merchant{Code: c.MerchantCode}
	if e = c.Merchant.Read("Code"); e != nil {
		o.Failure("merchant_code.invalid", util.ErrorInvalidData("merchant code"))
	} else {
		// Allow Customer to Login but they'll be validated by Merchant Suspended status
		if c.Merchant.Suspended == 1 {
			c.SuspendOrUnSuspend = "un-suspended"
			c.Merchant.Suspended = 2
		} else {
			c.SuspendOrUnSuspend = "suspended"
			c.Merchant.Suspended = 1
		}
	}

	return o
}

// Messages : function to return error messages after validation
func (c *suspensionRequest) Messages() map[string]string {
	return map[string]string{}
}
