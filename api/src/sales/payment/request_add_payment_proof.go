// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// addPaymentProofRequest : struct to hold request data for adding payment proof on sales payment
type addPaymentProofRequest struct {
	ID       int64  `json:"-"`
	ImageUrl string `json:"image_url" valid:"required"`

	SalesPayment *model.SalesPayment `json:"-"`
	Session      *auth.SessionData   `json:"-"`
}

// Validate : function to validate Add Payment Proof on Sales Payment request data
func (r *addPaymentProofRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	r.SalesPayment = &model.SalesPayment{ID: r.ID}
	if err = r.SalesPayment.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales payment"))
		return o
	}

	if r.ImageUrl == "" {
		o.Failure("image_url.required", util.ErrorInputRequired("image url"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *addPaymentProofRequest) Messages() map[string]string {
	messages := map[string]string{
		"image_url.required": util.ErrorInputRequired("image url"),
	}

	return messages
}
