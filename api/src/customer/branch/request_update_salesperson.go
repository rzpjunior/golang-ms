// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updatesalespersonRequest struct {
	ID            int64  `json:"-" valid:"required"`
	SalespersonID string `json:"salesperson_id" valid:"required"`

	Salesperson *model.Staff
	Branch      *model.Branch

	Session *auth.SessionData
}

// Validate : function to validate supplier request data
func (c *updatesalespersonRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	c.Branch = &model.Branch{ID: c.ID}
	if err = c.Branch.Read("ID"); err == nil {
		if c.Branch.Status != 1 {
			o.Failure("branch.inactive", util.ErrorActive("outlet"))
		} else {
			c.Branch.Salesperson.Read("ID")
		}
	} else {
		o.Failure("branch.invalid", util.ErrorInvalidData("outlet"))
	}

	SalesPersonID, _ := common.Decrypt(c.SalespersonID)
	c.Salesperson = &model.Staff{ID: SalesPersonID}
	if err = c.Salesperson.Read("ID"); err != nil {
		o.Failure("salesperson.invalid", util.ErrorInvalidData("salesperson"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *updatesalespersonRequest) Messages() map[string]string {
	return map[string]string{
		"salesperson_id.required": util.ErrorInputRequired("salesperson_id"),
	}
}

type updateBulkSalespersonReq struct {
	Data []*dataBulkUpdateSalesperson `json:"data" valid:"required"`

	Session *auth.SessionData
}

type dataBulkUpdateSalesperson struct {
	BranchCode         string `json:"branch_code" valid:"required"`
	NewSalesPersonCode string `json:"new_salesperson_code" valid:"required"`

	Branch         *model.Branch `json:"-"`
	NewSalesPerson *model.Staff  `json:"-"`
}

// Validate : function to validate bulk update salesperson request data
func (c *updateBulkSalespersonReq) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	for i, v := range c.Data {
		if v.BranchCode != "" {
			v.Branch = &model.Branch{Code: v.BranchCode}
			err = v.Branch.Read("Code")
			if err != nil {
				o.Failure("branch_code_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("outlet"))
			}

			if v.Branch.ID != 0 {
				if v.Branch.Status != 1 {
					o.Failure("branch_code_"+strconv.Itoa(i)+".inactive", util.ErrorActive("outlet"))
				}

				if err = v.Branch.Salesperson.Read("ID"); err != nil {
					o.Failure("salesperson_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales person"))
				}
			}
		}

		if v.NewSalesPersonCode != "" {
			v.NewSalesPerson = &model.Staff{Code: v.NewSalesPersonCode}
			if err = v.NewSalesPerson.Read("Code"); err != nil {
				o.Failure("new_salesperson_code_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("new sales person code"))
			}
			if v.NewSalesPerson.Status != 1 {
				o.Failure("new_salesperson_code_"+strconv.Itoa(i)+".inactive", util.ErrorActive("new sales person code"))
			}

			if v.Branch.Salesperson.ID == v.NewSalesPerson.ID {
				o.Failure("new_salesperson_code_"+strconv.Itoa(i)+".invalid", util.ErrorMustDifferenctSalesPerson(v.BranchCode, v.NewSalesPersonCode))
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateBulkSalespersonReq) Messages() map[string]string {
	return map[string]string{
		"data.required":                 util.ErrorInputRequired("data"),
		"branch_code.required":          util.ErrorInputRequired("branch code"),
		"new_salesperson_code.required": util.ErrorInputRequired("new Salesperson code"),
	}
}
