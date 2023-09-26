// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type resetPasswordRequest struct {
	ID              int64  `json:"-" valid:"required"`
	Password        string `json:"password" valid:"required"`
	ConfirmPassword string `json:"confirm_password" valid:"required"`

	PasswordHash string            `json:"-"`
	Session      *auth.SessionData `json:"-"`
}

func (c *resetPasswordRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	if errors := util.CheckPassword(c.Password); errors != "" {
		o.Failure("password.invalid", errors)
	}
	//validation password
	if c.ConfirmPassword != c.Password {
		o.Failure("confirm_password.notmatch", "password not match")
	}
	if c.PasswordHash, err = common.PasswordHasher(c.Password); err != nil {
		o.Failure("password.invalid", util.ErrorInvalidData("password"))
	}

	return o
}

func (c *resetPasswordRequest) Messages() map[string]string {
	return map[string]string{
		"password.required":         util.ErrorInputRequired("password"),
		"confirm_password.required": util.ErrorInputRequired("confirm password"),
	}
}
