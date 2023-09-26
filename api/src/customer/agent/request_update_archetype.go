// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateArchetypeRequest struct {
	ID             int64  `json:"-" valid:"required"`
	ArchetypeId    string `json:"archetype_id" valid:"required"`
	PrevArchetype  string `json:"prev_archetype"`
	BusinessTypeID string `json:"business_type"`

	Archetype    *model.Archetype
	BusinessType *model.BusinessType

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateArchetypeRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	//var err error

	if archetypeId, err := common.Decrypt(c.ArchetypeId); err == nil {
		if c.Archetype, err = repository.ValidArchetype(archetypeId); err != nil {
			o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
		} else {
			if c.Archetype.Status != int8(1) {
				o.Failure("archetype_id.active", util.ErrorActive("archetype"))
			} else {
				if businessTypeId, err := common.Decrypt(c.BusinessTypeID); err == nil {
					if c.BusinessType, err = repository.ValidBusinessType(businessTypeId); err != nil {
						o.Failure("business_type.invalid", util.ErrorInvalidData("business type"))
					}
				}
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateArchetypeRequest) Messages() map[string]string {
	return map[string]string{
		"archetype.required": util.ErrorInputRequired("archetype"),
	}
}
