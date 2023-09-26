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

type archiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var filter, exclude map[string]interface{}
	var total int64

	if archetype, err := repository.ValidArchetype(c.ID); err == nil {
		if archetype.Status != 1 {
			o.Failure("id.invalid", util.ErrorActive("status"))
		}

		filter = map[string]interface{}{"status__in": []int{1, 2}, "archetype_id": c.ID}
		if _, total, err = repository.CheckBranchData(filter, exclude); err == nil && total > 0 {
			o.Failure("id.invalid", util.ErrorRelated("active or archive ", "customer", "archetype"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("archetype"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
