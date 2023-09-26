// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updatePhoneRequest struct {
	ID          int64  `json:"-" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"required"`

	Merchant *model.Merchant

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate merchant request data
func (c *updatePhoneRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var existPhoneNum bool
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}

		c.Merchant.UserMerchant.Read("ID")
	} else {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	existPhoneNum = orSelect.QueryTable("merchant").Filter("phone_number", strings.TrimPrefix(c.PhoneNumber, "0")).Filter("status__in", 1, 2).Exist()
	if existPhoneNum {
		o.Failure("phone_number.invalid", util.ErrorPhoneNumber())
	}

	return o
}

// Messages : function to return error validation messages
func (c *updatePhoneRequest) Messages() map[string]string {
	return map[string]string{
		"phone_number.required": util.ErrorInputRequired("phone number"),
	}
}
