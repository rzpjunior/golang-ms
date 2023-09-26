// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Merchant *model.Merchant `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 2 {
			o.Failure("status.archive", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("agent.invalid", util.ErrorInvalidData("agent"))
	}

	return o
}

func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
