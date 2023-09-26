// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_group

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID      int64             `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`

	SalesGroup *model.SalesGroup
	Staff      []*model.Staff
}

// Validate : function to validate request data
func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q := orm.NewOrm()
	q.Using("read_only")
	var err error

	if c.SalesGroup, err = repository.ValidSalesGroup(c.ID); err != nil {
		o.Failure("sales_group.invalid", util.ErrorInvalidData("sales_group"))
	}

	if c.SalesGroup.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("status"))
	}

	// get staff related to sales group
	q.Raw("SELECT * FROM staff where sales_group_id = ?", c.SalesGroup.ID).QueryRows(&c.Staff)

	return o
}

// Messages : function to return error messages after validation
func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
