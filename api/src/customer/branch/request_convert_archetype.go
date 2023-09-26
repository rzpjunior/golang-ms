// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type convertarchetypeRequest struct {
	ID          int64  `json:"-" valid:"required"`
	ArchetypeID string `json:"archetype_id" valid:"required"`

	Archetype *model.Archetype
	Branch    *model.Branch

	Session *auth.SessionData
}

// Validate : function to validate supplier request data
func (c *convertarchetypeRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	archetypeID, _ := common.Decrypt(c.ArchetypeID)
	c.Archetype = &model.Archetype{ID: archetypeID}
	if err = c.Archetype.Read("ID"); err != nil {
		o.Failure("archetype.invalid", util.ErrorInvalidData("archetype"))
	}

	c.Branch = &model.Branch{ID: c.ID}
	if err = c.Branch.Read("ID"); err == nil {
		c.Branch.Merchant = &model.Merchant{ID: c.Branch.Merchant.ID}
		if err = c.Branch.Merchant.Read("ID"); err != nil {
			o.Failure("merchant.invalid", util.ErrorInvalidData("main outlet"))
		}

		err = c.Branch.Archetype.Read("ID")
		if err != nil {
			o.Failure("archetype.invalid", util.ErrorInvalidData("archetype"))
		}

		err = c.Archetype.BusinessType.Read("ID")
		if err != nil {
			o.Failure("business_type.invalid", util.ErrorInvalidData("business type"))
		}
		if c.Archetype.BusinessType.ID != c.Branch.Archetype.BusinessType.ID {
			o.Failure("archetype_id.invalid", util.ErrorCannotCrossBusinessType())
		}
	} else {
		o.Failure("branch.invalid", util.ErrorInvalidData("outlet"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *convertarchetypeRequest) Messages() map[string]string {
	return map[string]string{
		"archetype_id.required": util.ErrorInputRequired("archetype_id"),
	}
}
