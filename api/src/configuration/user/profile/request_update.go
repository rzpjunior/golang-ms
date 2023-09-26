// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package profile

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	DisplayName string ` json:"display_name" valid:"required"`
	PhoneNumber string ` json:"phone_number" valid:"required"`

	Staff   *model.Staff
	Session *auth.SessionData `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	c.Staff = c.Session.Staff
	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"display_name.required": util.ErrorInputRequired("display name"),
		"phone_number.required": util.ErrorInputRequired("phone number"),
	}
}
