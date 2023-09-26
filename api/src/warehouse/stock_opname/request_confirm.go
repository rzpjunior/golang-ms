// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock_opname

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type confirmRequest struct {
	ID          int64              `json:"-"`
	Session     *auth.SessionData  `json:"-"`
	StockOpname *model.StockOpname `json:"-"`
}

// Validate : function to validate uom request data
func (c *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	c.StockOpname = &model.StockOpname{ID: c.ID}
	if e = c.StockOpname.Read("ID"); e == nil {
		if c.StockOpname.Status != 1 {
			o.Failure("id.invalid", util.ErrorActive("stock opname"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("stock opname"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
