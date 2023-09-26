// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strconv"
)

// updateRequestAssign : struct to hold picking assign request data
type updateRequestAssign struct {
	ID int64 `json:"-"`

	PickingOrderItems  []*pickingItems           `json:"picking_order_items" valid:"required"`
	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	TypeRequest        string                    `json:"type_request"`

	Session *auth.SessionData `json:"-"`
}

type pickingItems struct {
	PickingOrderItemID string  `json:"picking_order_item_id"`
	PickOrderQty       float64 `json:"pick_qty"`
	CheckOrderQty      float64 `json:"check_qty"`

	PickingOrderItem *model.PickingOrderItem `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *updateRequestAssign) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status == int8(2) {
			o.Failure("picking_order_assign.active", util.ErrorPickingStatus("new", "on progress"))
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

	for i, v := range r.PickingOrderItems {
		if pickingOrderAssign, err := common.Decrypt(v.PickingOrderItemID); err == nil {
			if v.PickingOrderItem, err = repository.ValidPickingOrderItem(pickingOrderAssign); err == nil {
				if v.PickOrderQty < 0 {
					o.Failure("id.invalid", util.ErrorGreater("pick quantity", "0"))
					return o
				}
			} else {
				o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
			}
		} else {
			o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *updateRequestAssign) Messages() map[string]string {
	messages := map[string]string{
		"picking_order_items.required": util.ErrorInputRequired("picking order items"),
	}

	return messages
}
