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

// countPrintRequest : struct to hold count number of print copy request data
type countPrintRequest struct {
	ID int64 `json:"-" valid:"required"`

	PurchaseOrder *model.PurchaseOrder `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

// Validate : function to validate print request data
func (r *countPrintRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	r.PurchaseOrder, err = repository.ValidPurchaseOrder(r.ID)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order id"))
	}

	if r.PurchaseOrder.Status == 3 {
		o.Failure("purchase_order.invalid", util.ErrorStatusDoc("purchase order", "printed", "purchase order"))
	}

	return o
}

func (r *countPrintRequest) Messages() map[string]string {
	return map[string]string{}
}
