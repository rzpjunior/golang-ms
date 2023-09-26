// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package archetype

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold request data
type createRequest struct {
	Code           string `json:"-"`
	Name           string `json:"name" valid:"required"`
	BusinessTypeId string `json:"business_type_id" valid:"required"`
	Note           string `json:"note"`
	CustomerGroup  string `json:"customer_group" valid:"required"`
	// Abbreviation   string `json:"abbreviation"`

	BusinessType *model.BusinessType

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("archetype"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if businessTypeId, e := common.Decrypt(c.BusinessTypeId); e != nil {
		o.Failure("business_type.invalid", util.ErrorInvalidData("business type"))
	} else {
		if c.BusinessType, e = repository.ValidBusinessType(businessTypeId); e != nil {
			o.Failure("business_type.invalid", util.ErrorInvalidData("business type"))
		} else {
			if c.BusinessType.Status != int8(1) {
				o.Failure("business_type.invalid", util.ErrorActive("business type"))
			}
		}
	}

	filter := map[string]interface{}{"name": c.Name}
	exclude := map[string]interface{}{"status": 3}
	if _, countArchetype, err := repository.CheckArchetypeData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countArchetype > 0 {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	// abbreviation put on hold for the moment
	// archetypeAbv := &model.Archetype{Abbreviation: c.Abbreviation}
	// if err = archetypeAbv.Read("Abbreviation"); err == nil {
	// 	o.Failure("abbreviation", util.ErrorDuplicate("abbreviation"))
	// }

	return o
}

// Messages : function to return error validation messages after validation
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":             util.ErrorInputRequired("name"),
		"business_type_id.required": util.ErrorInputRequired("business type"),
		"customer_group.required":   util.ErrorInputRequired("customer group"),
		// "abbreviation.required": util.ErrorInputRequired("abbreviation"),
	}
}
