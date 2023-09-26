// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequestApprovalAssign : struct to hold picking assign request data
type updateRequestApprovalAssign struct {
	ID int64 `json:"-"`

	TotalColly         float64                   `json:"total_colly" valid:"required"`
	PickingOrderItems  []*approvalPickingItems   `json:"picking_order_items" valid:"required"`
	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	DeliveryKolies     []*DeliveryKoli           `json:"delivery_kolies" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type approvalPickingItems struct {
	PickingOrderItemID string  `json:"picking_order_item_id"`
	PickOrderQty       float64 `json:"pick_qty"`
	CheckOrderQty      float64 `json:"check_qty"`
	UnfullfillNote     string  `json:"unfullfill_note"`
	POIFlagging        int8    `json:"poi_flagging"`

	PickingOrderItem *model.PickingOrderItem `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *updateRequestApprovalAssign) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.TotalColly <= 0 {
		o.Failure("id.invalid", util.ErrorEqualGreaterInd("total koli", "0"))
		return o
	}

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status == int8(5) || r.PickingOrderAssign.Status == int8(4) {
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

	for i2, v2 := range r.DeliveryKolies {
		if koliID, err := common.Decrypt(v2.KoliID); err == nil {
			if v2.Koli, err = repository.ValidKoli(koliID); err != nil {
				o.Failure("koli_id"+strconv.Itoa(i2)+".invalid", util.ErrorInvalidData("koli"))
			}
		} else {
			o.Failure("koli_id"+strconv.Itoa(i2)+".invalid", util.ErrorInvalidData("koli"))
		}
		v2.SalesOrder = r.PickingOrderAssign.SalesOrder
	}

	return o
}

// Messages : function to return error validation messages
func (r *updateRequestApprovalAssign) Messages() map[string]string {
	messages := map[string]string{
		"picking_order_items.required": util.ErrorInputRequired("picking order items"),
	}

	return messages
}
