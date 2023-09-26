// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	CodeRole     string              `json:"-"`
	Name         string              ` json:"name" valid:"required"`
	DivisionID   string              ` json:"division_id" valid:"required"`
	Note         string              `json:"note"`
	PermissionID []string            `json:"permission_id" valid:"required"`
	Permission   []*model.Permission `json:"-"`

	Session  *auth.SessionData `json:"-"`
	Division *model.Division   `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var id int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.CodeRole, err = util.CheckTable("role"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code user"))
	}

	// validasi name tidak boleh duplikat
	orSelect.Raw("SELECT COUNT(id) FROM role WHERE name = ?  AND status != 3", c.Name).QueryRow(&id)
	if id > 0 {
		o.Failure("name.invalid", util.ErrorDuplicate("name"))
	}

	if c.DivisionID != "" {
		if divID, e := common.Decrypt(c.DivisionID); e != nil {
			o.Failure("division_id.invalid", util.ErrorInvalidData("division"))
		} else {
			if c.Division, e = repository.ValidDivision(divID); e != nil {
				o.Failure("division_id.invalid", util.ErrorInvalidData("division"))
			}
		}
	}
	if len(c.PermissionID) < 1 {
		o.Failure("id.invalid", util.ErrorSelectOne("permission"))
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

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":        util.ErrorInputRequired("name"),
		"division_id.required": util.ErrorSelectRequired("division"),
	}
}
