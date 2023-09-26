// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// checkinBulkRequestAssign : struct to hold picking assign request data
type checkinBulkRequestAssign struct {
	HelperID    string       `json:"staff_id" valid:"required"`
	SalesOrder  []salesOrder `json:"sales_order"`
	TypeRequest string       `json:"type_request"`

	PickingOrderAssign []*model.PickingOrderAssign `json:"-"`

	IsFinished   bool                `json:"-"`
	PickingOrder *model.PickingOrder `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type salesOrder struct {
	SalesOrderID string `json:"sales_order_id"`
}

// Validate : function to validate picking assign request data
func (r *checkinBulkRequestAssign) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error

	helperID, _ := common.Decrypt(r.HelperID)
	var soID int64
	for i, v := range r.SalesOrder {
		if soID, err = common.Decrypt(v.SalesOrderID); err != nil {
			o.Failure("sales_order.invalid", util.ErrorInvalidData("sales order"))
		}

		var salesOrder *model.SalesOrder
		var pickingOrderAssign *model.PickingOrderAssign
		var isExists int
		if salesOrder, err = repository.ValidSalesOrder(soID); err != nil {
			o.Failure("sales_order.invalid", util.ErrorInvalidData("sales order"))
		}

		o1.Raw("select EXISTS (SELECT so.status from sales_order so where so.id = ? and so.status in (1,3,9,12))", salesOrder.ID).QueryRow(&isExists)
		if isExists != 1 {
			o.Failure("status.invalid", util.ErrorSalesOrderOnPicking())
		}

		o1.Raw("SELECT id, picking_order_id, sales_order_id, staff_id, courier_id, courier_vendor_id, dispatcher_id, picking_list_id, status, dispatch_status, dispatch_timestamp, assign_timestamp, planning_vendor, been_rejected, note, checkin_timestamp, checkout_timestamp, checker_in_timestamp, checker_out_timestamp, total_koli, total_scan_dispatch, checked_at, checked_by "+
			"from eden_v2.picking_order_assign where sales_order_id = ?", salesOrder.ID).QueryRow(&pickingOrderAssign)

		if pickingOrderAssign.Helper.ID != helperID {
			o.Failure("helper.invalid", util.ErrorInvalidData("helper"))
		}

		r.PickingOrderAssign = append(r.PickingOrderAssign, pickingOrderAssign)

		// read to db just only in last element
		if i == len(r.SalesOrder)-1 {
			pickingOrderAssign.PickingOrder.Read("ID")
			r.PickingOrder = &model.PickingOrder{
				ID: pickingOrderAssign.PickingOrder.ID,
			}
		}

	}

	return o
}

// Messages : function to return error validation messages
func (r *checkinBulkRequestAssign) Messages() map[string]string {
	messages := map[string]string{
		"staff_id.required": util.ErrorInputRequired("helper"),
	}

	return messages
}
