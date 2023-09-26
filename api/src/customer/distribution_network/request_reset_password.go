// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package distribution_network

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// resetPasswordRequest data struct that store request data when requesting an auth creation.
type resetPasswordRequest struct {
	ID int64 `json:"-" valid:"required"`

	PasswordHash string `json:"-"`

	Merchant *model.Merchant `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *resetPasswordRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err             error
		o1              = orm.NewOrm()
		defaultPassword = "12345678"
	)
	o1.Using("read_only")

	// Check Merchant
	if r.Merchant, err = repository.ValidMerchant(r.ID); err != nil {
		o.Failure("merchant_id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	// Check if merchant belongs to EDN
	if err = r.Merchant.BusinessType.Read("ID"); err != nil {
		o.Failure("business_type.invalid", util.ErrorInvalidData("business type"))
		return o
	}

	if r.Merchant.BusinessType.Name != "EDN" {
		o.Failure("business_type.invalid", util.ErrorMustBeSame("business type", "EDN"))
		return o
	}

	// Check if password is match with body request
	if r.PasswordHash, err = common.PasswordHasher(defaultPassword); err != nil {
		o.Failure("password.invalid", util.ErrorInvalidData("password"))
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *resetPasswordRequest) Messages() map[string]string {
	return map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}
}
