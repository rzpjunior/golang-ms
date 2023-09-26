// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateHelperRequest struct {
	ID   int64  `json:"-" valid:"required"`
	Note string `json:"note"`
	// staff
	Name            string ` json:"name" valid:"required"`
	RoleID          string ` json:"role_id" valid:"required"`
	WarehouseID     string ` json:"warehouse_id" valid:"required"`
	PhoneNumber     string ` json:"phone_number" valid:"required"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	PasswordHash    string `json:"-"`

	Role      *model.Role
	Staff     *model.Staff
	Parent    *model.Staff
	Session   *auth.SessionData
	Warehouse *model.Warehouse

	OldStaff *model.Staff
	OldUser  *model.User
}

func (c *updateHelperRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	if c.RoleID != "" {
		if roleID, e := common.Decrypt(c.RoleID); e != nil {
			o.Failure("role_id.invalid", util.ErrorInvalidData("role"))
		} else {
			if c.Role, e = repository.ValidRole(roleID); e != nil {
				o.Failure("role_id.invalid", util.ErrorInvalidData("role"))
			}
		}
	}

	if c.WarehouseID != "" {
		if warehouseID, e := common.Decrypt(c.WarehouseID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		} else {
			if c.Warehouse, e = repository.ValidWarehouse(warehouseID); e != nil {
				o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
			}
		}
	}

	c.Staff = &model.Staff{ID: c.ID}
	if err = c.Staff.Read("ID"); err != nil {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	c.OldStaff = c.Staff

	if err = c.Staff.User.Read("ID"); err != nil {
		o.Failure("user.invalid", util.ErrorInvalidData("user"))
	}

	c.OldUser = c.Staff.User

	if c.Password != "" {
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
	}

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number", util.ErrorCharLength("phone number", 8))
	}

	return o
}

func (c *updateHelperRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":         util.ErrorInputRequired("name"),
		"role_id.required":      util.ErrorSelectRequired("role"),
		"warehouse_id.required": util.ErrorSelectRequired("warehouse"),
		"phone_number.required": util.ErrorInputRequired("phone number"),
	}
}
