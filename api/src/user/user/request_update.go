// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateRequest struct {
	ID   int64  `json:"-" valid:"required"`
	Note string `json:"note"`
	// staff
	Name          string ` json:"name" valid:"required"`
	DisplayName   string ` json:"display_name" valid:"required"`
	DivisionID    string ` json:"division_id" valid:"required"`
	RoleID        string ` json:"role_id" valid:"required"`
	WarehouseID   string ` json:"warehouse_id" valid:"required"`
	AreaID        string ` json:"area_id" valid:"required"`
	SuperVisorID  string ` json:"supervisor_id"`
	PhoneNumber   string ` json:"phone_number" valid:"required"`
	SalesGroupID  string ` json:"sales_group_id"`
	SalesGroupInt int64
	RoleGroup     int8

	OldStaff   *model.Staff
	OldUser    *model.User
	Role       *model.Role
	Area       *model.Area
	Staff      *model.Staff
	Parent     *model.Staff
	Session    *auth.SessionData
	Warehouse  *model.Warehouse
	SalesGroup *model.SalesGroup
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	user := &model.User{ID: c.ID}
	if err = user.Read("ID"); err != nil {
		o.Failure("user.invalid", util.ErrorInvalidData("user"))
	}

	c.OldUser = user

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
				var attribute string
				orSelect.Raw("SELECT id, attribute FROM config_app WHERE attribute like '%%role_id' AND concat(',', value, ',') LIKE ? ", "%,"+strconv.FormatInt(c.Role.ID, 10)+",%").QueryRow(&id, &attribute)
				if id > 0 {
					if strings.Contains(attribute, "salesperson") {
						c.RoleGroup = 1
					} else {
						c.RoleGroup = 2
					}
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

	c.Staff = &model.Staff{User: user}
	if err = c.Staff.Read("User"); err != nil {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	c.OldStaff = c.Staff

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

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":         util.ErrorInputRequired("name"),
		"display_name.required": util.ErrorInputRequired("display name"),
		"role_id.required":      util.ErrorSelectRequired("role"),
		"area_id.required":      util.ErrorSelectRequired("area"),
		"division_id.required":  util.ErrorSelectRequired("division"),
		"warehouse_id.required": util.ErrorSelectRequired("warehouse"),
		"phone_number.required": util.ErrorInputRequired("phone number"),
	}
}
