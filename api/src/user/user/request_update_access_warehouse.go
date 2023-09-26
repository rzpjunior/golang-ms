// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

type updateWarehouseAccessRequest struct {
	ID               int64    `json:"-" valid:"required"`
	WarehouseStr     string   `json:"-"`
	WarehouseChecked []string `json:"warehouse_checked"`

	User    *model.User       `json:"user"`
	Staff   *model.Staff      `json:"staff"`
	Session *auth.SessionData `json:"-"`
}

func (c *updateWarehouseAccessRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	user := &model.User{ID: c.ID}
	if err = user.Read("ID"); err != nil {
		o.Failure("user.invalid", util.ErrorInvalidData("user"))
	}

	c.Staff = &model.Staff{User: user}
	if err = c.Staff.Read("User"); err != nil {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	for i, v := range c.WarehouseChecked {
		v = common.Encrypt(v)
		c.WarehouseChecked[i] = v
		c.WarehouseStr = c.WarehouseStr + v + ","
	}
	c.WarehouseStr = strings.TrimSuffix(c.WarehouseStr, ",")

	return o
}

func (c *updateWarehouseAccessRequest) Messages() map[string]string {
	return map[string]string{}
}
