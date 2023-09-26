// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateParamRequest struct {
	ID              int64   `json:"-" valid:"required"`
	LimitSalesOrder int8    `json:"limit_sales_order" valid:"required"`
	LimitWeight     float64 `json:"limit_weight" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateParamRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (c *updateParamRequest) Messages() map[string]string {
	return map[string]string{
		"limit_sales_order.required": util.ErrorInputRequired("limit sales order"),
		"limit_weight.required":      util.ErrorInputRequired("limit weight"),
	}
}
