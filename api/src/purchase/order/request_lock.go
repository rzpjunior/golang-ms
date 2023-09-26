// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// lockRequest : struct to lock purchase order
type lockRequest struct {
	ID        int64 `json:"-" valid:"required"`
	CancelReq int8  `json:"-"`

	PurchaseOrder *model.PurchaseOrder `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate purchase order request data
func (r *lockRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.PurchaseOrder, err = repository.ValidPurchaseOrder(r.ID); err != nil {
		o.Failure("purchase_order.invalid", util.ErrorInvalidData("purchase order"))
		return o
	}

	// To unlock when press cancel in
	if r.PurchaseOrder.Locked == 1 {
		r.CancelReq = 1
	}

	return o
}

// Messages : function to return error validation messages
func (r *lockRequest) Messages() map[string]string {
	messages := map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}

	return messages
}
