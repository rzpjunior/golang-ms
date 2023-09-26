// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// cancelRequest : struct to hold price set request data
type cancelActiveRequest struct {
	ID   int64  `json:"-"`
	Note string `json:"note" valid:"required"`

	SalesPayment *model.SalesPayment `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (r *cancelActiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	r.SalesPayment = &model.SalesPayment{ID: r.ID}
	if err = r.SalesPayment.Read("ID"); err == nil {
		if r.SalesPayment.Status != 1 {
			o.Failure("status.inactive", util.ErrorDocStatus("sales payment", "active"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("sales payment"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelActiveRequest) Messages() map[string]string {
	messages := map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}

	return messages
}
