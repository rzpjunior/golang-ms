// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package archetype

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

// Validate : function to validate request data
func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if archetype, err := repository.ValidArchetype(c.ID); err == nil {
		if archetype.Status != 2 {
			o.Failure("id.invalid", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("archetype"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
