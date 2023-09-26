// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updatePhoneNumber struct {
	ID          int64  `json:"-" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"required"`

	Merchant     *model.Merchant
	UserMerchant *model.UserMerchant
	Session      *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updatePhoneNumber) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var existPhoneNum bool
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
		if c.UserMerchant, err = repository.ValidUserMerchant(c.Merchant.UserMerchant.ID); err != nil {
			o.Failure("user_merchant.invalid", util.ErrorInvalidData("user merchant"))
		}
	} else {
		o.Failure("agent.invalid", util.ErrorInvalidData("agent"))
	}

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number.invalid", util.ErrorCharLength("phone number", 8))
	}

	existPhoneNum = orSelect.QueryTable("merchant").Filter("phone_number", strings.TrimPrefix(c.PhoneNumber, "0")).Filter("status__in", 1, 2).Exist()
	if existPhoneNum {
		o.Failure("phone_number.invalid", util.ErrorPhoneNumber())
	}

	return o
}

// Messages : function to return error validation messages
func (c *updatePhoneNumber) Messages() map[string]string {
	return map[string]string{
		"phone_number.required": util.ErrorInputRequired("phone number"),
	}
}
