// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// checkinRequestChecker : struct to hold picking assign request data
type checkinRequestChecker struct {
	ID int64 `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	SalesOrderItems    []*model.SalesOrderItem   `json:"-"`

	IsFinished bool `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *checkinRequestChecker) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	o1 := orm.NewOrm()
	o1.Using("read_only")

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status != int8(5) {
			o.Failure("picking_order_assign.active", util.ErrorPickingSingleStatus("picked"))
		}
		r.PickingOrderAssign.SalesOrder.Read("ID")
		if r.PickingOrderAssign.SalesOrder.Status != 1 &&
			r.PickingOrderAssign.SalesOrder.Status != 9 &&
			r.PickingOrderAssign.SalesOrder.Status != 12 {
			o.Failure("status.invalid", util.ErrorSalesOrderOnPicking())
			return o
		}
		if r.PickingOrderAssign.SalesOrder.IsLocked == 1 {
			o.Failure("sales_order.invalid", util.ErrorSOLockedInd())
			return o
		}
		r.PickingOrderAssign.PickingOrder.Read("ID")
		o1.Raw("SELECT id, picking_order_assign_id, product_id, order_qty, pick_qty, check_qty, unfullfill_note FROM eden_v2.picking_order_item where picking_order_assign_id = ?", r.PickingOrderAssign.ID).QueryRows(&r.PickingOrderAssign.PickingOrderItem)
		o1.Raw("SELECT id, sales_order_id, product_id, order_qty, forecast_qty, unit_price, shadow_price, subtotal, weight, note FROM eden_v2.sales_order_item where sales_order_id = ?", r.PickingOrderAssign.SalesOrder.ID).QueryRows(&r.SalesOrderItems)

	} else {
		o.Failure("picking_order_assign.invalid", util.ErrorInvalidData("picking order assign"))
	}

	return o
}

// Messages : function to return error validation messages
func (r *checkinRequestChecker) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
