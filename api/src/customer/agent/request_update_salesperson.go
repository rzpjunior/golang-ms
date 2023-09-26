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

type updateSalespersonRequest struct {
	ID              int64  `json:"-" valid:"required"`
	SalespersonId   string `json:"salesperson_id" valid:"required"`
	PrevSalesperson string `json:"prev_salesperson"`

	Salesperson *model.Staff

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateSalespersonRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if salesPersonId, err := common.Decrypt(c.SalespersonId); err != nil {
		o.Failure("salesperson.invalid", util.ErrorInvalidData("salesperson"))
	} else {
		if c.Salesperson, err = repository.ValidStaff(salesPersonId); err != nil {
			o.Failure("salesperson.invalid", util.ErrorInvalidData("salesperson"))
		} else {
			if c.Salesperson.Status != int8(1) {
				o.Failure("salesperson.active", util.ErrorActive("salesperson"))
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateSalespersonRequest) Messages() map[string]string {
	return map[string]string{
		"salesperson.required": util.ErrorInputRequired("salesperson"),
	}
}
