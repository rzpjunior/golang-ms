// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// checkinRequestAssign : struct to hold picking assign request data
type checkinRequestAssign struct {
	ID int64 `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`

	IsFinished bool `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *checkinRequestAssign) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status != int8(1) {
			o.Failure("picking_order_assign.active", util.ErrorPickingSingleStatus("new"))
		}
		r.PickingOrderAssign.SalesOrder.Read("ID")
		if r.PickingOrderAssign.SalesOrder.Status != 1 &&
			r.PickingOrderAssign.SalesOrder.Status != 9 &&
			r.PickingOrderAssign.SalesOrder.Status != 12 {
			o.Failure("status.invalid", util.ErrorSalesOrderOnPicking())
		}
		r.PickingOrderAssign.PickingOrder.Read("ID")
	} else {
		o.Failure("picking_order_assign.invalid", util.ErrorInvalidData("picking order assign"))
	}

	return o
}

// Messages : function to return error validation messages
func (r *checkinRequestAssign) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
