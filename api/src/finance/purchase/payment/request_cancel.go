// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_payment

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type cancelRequest struct {
	ID               int64  `json:"-"`
	CancellationNote string `json:"note" valid:"required"`
	PurchasePayment  *model.PurchasePayment

	Session *auth.SessionData
}

// Validate : function to validate request data
func (c *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if len(c.CancellationNote) > 250 {
		o.Failure("id.invalid", util.ErrorCharLength("note", 250))
	}
	if c.PurchasePayment.Status != 2 {
		o.Failure("id.invalid", util.ErrorDocStatus("purchase payment", "finished"))
	}
	return o
}

// Messages : function to return error messages after validation
func (c *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
