// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createHelperRequest struct {
	CodeUser        string `json:"-"`
	Email           string `json:"email" valid:"required|email"`
	Password        string `json:"password" valid:"required"`
	ConfirmPassword string `json:"confirm_password" valid:"required"`
	PasswordHash    string `json:"-"`
	// staff
	Name        string ` json:"name" valid:"required"`
	RoleID      string ` json:"role_id" valid:"required"`
	WarehouseID string ` json:"warehouse_id" valid:"required"`
	PhoneNumber string ` json:"phone_number" valid:"required"`

	CodeStaff string            `json:"-"`
	RoleGroup int8              `json:"-"`
	Role      *model.Role       `json:"-"`
	Area      *model.Area       `json:"-"`
	Parent    *model.Staff      `json:"-"`
	Session   *auth.SessionData `json:"-"`
	Warehouse *model.Warehouse  `json:"-"`
}

func (c *createHelperRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.CodeUser, err = util.CheckTable("user"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code user"))
	}
	if c.CodeStaff, err = util.CheckTable("staff"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code staff"))
	}
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
			} else {
				c.Area = c.Warehouse.Area
			}
		}
	}

	// validasi email tidak boleh duplikat
	user := &model.User{Email: c.Email}
	if err = user.Read("Email"); err == nil {
		o.Failure("email.invalid", util.ErrorDuplicate("email"))
	}

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

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number.invalid", util.ErrorCharLength("phone number", 8))
	}

	return o
}

func (c *createHelperRequest) Messages() map[string]string {
	return map[string]string{
		"password.required":         util.ErrorInputRequired("password"),
		"confirm_password.required": util.ErrorInputRequired("confirm password"),
		"name.required":             util.ErrorInputRequired("name"),
		"role_id.required":          util.ErrorSelectRequired("role"),
		"warehouse_id.required":     util.ErrorSelectRequired("warehouse"),
		"email.required":            util.ErrorInputRequired("email"),
		"phone_number.required":     util.ErrorInputRequired("phone number"),
	}
}
