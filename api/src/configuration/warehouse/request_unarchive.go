// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if warehouse, err := repository.ValidWarehouse(c.ID); err == nil {
		if warehouse.Status != 2 {
			o.Failure("status.archive", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
	}

	return o
}

func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
