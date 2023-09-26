// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID           int64    `json:"-"`
	Name         string   ` json:"name" valid:"required"`
	Note         string   `json:"note"`
	PermissionID []string `json:"permission_id" valid:"required"`

	NewPermission       []int64           `json:"-"`
	OldRolePermissionID []int64           `json:"-"`
	Session             *auth.SessionData `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var id int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	// validasi name tidak boleh duplikat
	orSelect.Raw("SELECT COUNT(id) FROM role WHERE name = ? AND id != ? AND status != 3", c.Name, c.ID).QueryRow(&id)
	if id > 0 {
		o.Failure("name.invalid", util.ErrorDuplicate("name"))
	}

	if len(c.PermissionID) < 1 {
		o.Failure("id.invalid", util.ErrorInputRequired("permission"))
	} else {
		for _, id := range c.PermissionID {
			// this func for decrypt
			idConv, _ := common.Decrypt(id)
			// for get data permission from id
			// add permission
			c.NewPermission = append(c.NewPermission, idConv)
		}
		orSelect.Raw("SELECT permission_id FROM role_permission WHERE role_id = ?", c.ID).QueryRows(&c.OldRolePermissionID)

	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{}
}
