// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// lockRequest : struct to lock sales order
type lockRequest struct {
	ID        int64 `json:"-"`
	CancelReq int8  `json:"-"`

	SalesOrder     *model.SalesOrder       `json:"-"`
	SalesOrderItem []*model.SalesOrderItem `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate sales order request data
func (r *lockRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.SalesOrder, err = repository.ValidSalesOrder(r.ID); err == nil {
		if r.SalesOrder.Status != 1 {
			o.Failure("sales_order.inactive", util.ErrorActive("sales order"))
		}

		// To unlock when press cancel in create SO
		if r.SalesOrder.IsLocked == 1 {
			r.CancelReq = 1
		}
	} else {
		o.Failure("sales_order.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (r *lockRequest) Messages() map[string]string {
	messages := map[string]string{
		"note.required": util.ErrorInputRequired("note"),
	}

	return messages
}
