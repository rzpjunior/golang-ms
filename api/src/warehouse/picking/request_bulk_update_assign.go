// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateBulkQtyRequest : struct to hold picking assign request data
type updateBulkQtyRequest struct {
	ID        int64    `json:"-"`
	ProductID string   `json:"product_id"`
	Items     []*items `json:"items" valid:"required"`

	Product *model.Product    `json:"-"`
	Session *auth.SessionData `json:"-"`

	PickingList *model.PickingList `json:"-"`
}

type items struct {
	SalesOrderID   string  `json:"sales_order_id"`
	PickOrderQty   float64 `json:"pick_qty"`
	UnfullfillNote string  `json:"unfullfill_note"`
	PoiFlagging    int8    `json:"-"`

	SalesOrder *model.SalesOrder `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *updateBulkQtyRequest) Validate() *validation.Output {
	o1 := orm.NewOrm()
	o1.Using("read_only")
	o := &validation.Output{Valid: true}
	var err error

	if r.PickingList, err = repository.ValidPickingList(r.ID); err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}

	productID, _ := common.Decrypt(r.ProductID)
	if r.Product, err = repository.ValidProduct(productID); err != nil {
		o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		return o
	}

	for i, v := range r.Items {
		salesOrderID, _ := common.Decrypt(v.SalesOrderID)

		if v.SalesOrder, err = repository.ValidSalesOrder(salesOrderID); err != nil {
			o.Failure("sales_order_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales order"))
		}

		if err = o1.Raw("SELECT id, picking_order_id, sales_order_id, staff_id, courier_id, courier_vendor_id, dispatcher_id, picking_list_id, status, dispatch_status, dispatch_timestamp, assign_timestamp, planning_vendor, been_rejected, note, checkin_timestamp, checkout_timestamp, checker_in_timestamp, checker_out_timestamp, total_koli, total_scan_dispatch, checked_at, checked_by "+
			"FROM eden_v2.picking_order_assign where sales_order_id = ?", v.SalesOrder.ID).QueryRow(&v.PickingOrderAssign); err != nil {
			o.Failure("sales_order_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales order"))
		}

		if v.PickingOrderAssign.PickingList.ID != r.PickingList.ID {
			o.Failure("picking_list"+strconv.Itoa(i)+".invalid", util.ErrorMustBeSame("sales order", "grouping picking list"))
		}

		if v.UnfullfillNote != "" {
			v.PoiFlagging = 3
		} else {
			v.PoiFlagging = 2
		}

	}

	return o
}

// Messages : function to return error validation messages
func (r *updateBulkQtyRequest) Messages() map[string]string {
	messages := map[string]string{
		"items.required": util.ErrorInputRequired("items"),
	}

	return messages
}
