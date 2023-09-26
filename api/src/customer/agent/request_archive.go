// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Merchant *model.Merchant `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("agent.invalid", util.ErrorInvalidData("agent"))
	}

	// validation toward active sales order
	var branchs []model.Branch
	var count int
	orSelect.Raw("select * from branch b where merchant_id = ?", c.ID).QueryRows(&branchs)
	for _, v := range branchs {
		orSelect.Raw("select count(*) from sales_order so where so.branch_id =? and status = 1", v.ID).QueryRow(&count)
		if count > 0 {
			o.Failure("agent.invalid", util.ErrorRelated("active and archive ", "sales order", "agent"))
		}
	}

	return o
}

func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
