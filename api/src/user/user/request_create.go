// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	CodeUser        string              `json:"-"`
	Email           string              `json:"email" valid:"required|email"`
	Password        string              `json:"password" valid:"required"`
	ConfirmPassword string              `json:"confirm_password" valid:"required"`
	Note            string              `json:"note"`
	PermissionID    []string            `json:"permission_id" valid:"required"`
	PasswordHash    string              `json:"-"`
	Permission      []*model.Permission `json:"-"`
	// staff
	Name          string ` json:"name" valid:"required"`
	DisplayName   string ` json:"display_name" valid:"required"`
	EmployeeCode  string ` json:"employee_code" valid:"required"`
	DivisionID    string ` json:"division_id" valid:"required"`
	RoleID        string ` json:"role_id" valid:"required"`
	WarehouseID   string ` json:"warehouse_id" valid:"required"`
	AreaID        string ` json:"area_id" valid:"required"`
	SuperVisorID  string ` json:"supervisor_id"`
	PhoneNumber   string ` json:"phone_number" valid:"required"`
	SalesGroupID  string ` json:"sales_group_id"`
	SalesGroupInt int64

	CodeStaff  string            `json:"-"`
	RoleGroup  int8              `json:"-"`
	Role       *model.Role       `json:"-"`
	Area       *model.Area       `json:"-"`
	Parent     *model.Staff      `json:"-"`
	Session    *auth.SessionData `json:"-"`
	Warehouse  *model.Warehouse  `json:"-"`
	SalesGroup *model.SalesGroup `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.CodeUser, err = util.CheckTable("user"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code user"))
	}
	if c.CodeStaff, err = util.CheckTable("staff"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code staff"))
	}

	if areaID, e := common.Decrypt(c.AreaID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	} else {
		if c.Area, e = repository.ValidArea(areaID); e != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		} else {
			if c.Area.Status != int8(1) {
				o.Failure("area_id.invalid", util.ErrorActive("area"))
			}
		}
	}
	if c.RoleID != "" {
		if roleID, e := common.Decrypt(c.RoleID); e != nil {
			o.Failure("role_id.invalid", util.ErrorInvalidData("role"))
		} else {
			if c.Role, e = repository.ValidRole(roleID); e != nil {
				o.Failure("role_id.invalid", util.ErrorInvalidData("role"))
			} else {
				var id int
				orSelect.Raw("SELECT COUNT(id) FROM config_app WHERE attribute = 'salesperson_role_id' AND value LIKE ? ", "%"+strconv.FormatInt(c.Role.ID, 10)+"%").QueryRow(&id)
				if id == 1 {
					c.RoleGroup = 1
				} else {
					c.RoleGroup = 0
				}
			}
		}
	}

	if c.SuperVisorID != "" {
		if parentID, e := common.Decrypt(c.SuperVisorID); e != nil {
			o.Failure("parent_id.invalid", util.ErrorInvalidData("parent"))
		} else {
			if c.Parent, e = repository.ValidStaff(parentID); e != nil {
				o.Failure("parent_id.invalid", util.ErrorInvalidData("parent"))
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

	if len(c.PermissionID) < 1 {
		o.Failure("id.invalid", util.ErrorInputRequired("permission"))
	} else {
		for _, id := range c.PermissionID {
			// this func for decrypt
			idConv, _ := common.Decrypt(id)
			// for get data permission from id
			p := &model.Permission{ID: idConv}
			if err = p.Read("ID"); err == nil {
				// add permission
				c.Permission = append(c.Permission, p)
			}
		}
	}
	// this function for got id from array permission string

	// check sales group
	c.SalesGroupInt = 0
	if c.SalesGroupID != "" {
		if salesGroupID, e := common.Decrypt(c.SalesGroupID); e != nil {
			o.Failure("sales_group_id.invalid", util.ErrorInvalidData("sales group"))
		} else {
			if c.SalesGroup, e = repository.ValidSalesGroup(salesGroupID); e != nil {
				o.Failure("sales_group_id.invalid", util.ErrorInvalidData("sales group"))
			}
		}

		// check if sales group is active
		if c.SalesGroup.Status != 1 {
			o.Failure("sales_group_id.invalid", util.ErrorActive("sales group"))
		}
		c.SalesGroupInt = c.SalesGroup.ID
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"password.required":         util.ErrorInputRequired("password"),
		"confirm_password.required": util.ErrorInputRequired("confirm password"),
		"name.required":             util.ErrorInputRequired("name"),
		"display_name.required":     util.ErrorInputRequired("display name"),
		"role_id.required":          util.ErrorSelectRequired("role"),
		"area_id.required":          util.ErrorSelectRequired("area"),
		"division_id.required":      util.ErrorSelectRequired("division"),
		"warehouse_id.required":     util.ErrorSelectRequired("warehouse"),
		"email.required":            util.ErrorInputRequired("email"),
		"employee_code.required":    util.ErrorInputRequired("employee code"),
		"phone_number.required":     util.ErrorInputRequired("phone number"),
	}
}
