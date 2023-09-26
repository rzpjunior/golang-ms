// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updatePermissionRequest struct {
	ID           int64    `json:"-" valid:"required"`
	PermissionID []string `json:"permission_id" valid:"required"`

	NewPermission       []int64           `json:"-"`
	OldUserPermissionID []int64           `json:"-"`
	Session             *auth.SessionData `json:"-"`
}

func (c *updatePermissionRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if len(c.PermissionID) < 1 {
		o.Failure("id.invalid", util.ErrorInputRequired("permission"))
	}
	// this function for got id from array permission string
	for _, id := range c.PermissionID {
		// this func for decrypt
		idConv, _ := common.Decrypt(id)
		// for get data permission from id
		// add permission
		c.NewPermission = append(c.NewPermission, idConv)
	}
	orSelect.Raw("SELECT permission_id FROM user_permission WHERE user_id = ?", c.ID).QueryRows(&c.OldUserPermissionID)

	return o
}

func (c *updatePermissionRequest) Messages() map[string]string {
	return map[string]string{
		"permission_id.required": util.ErrorInputRequired("permission"),
	}
}
