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

// rejectRequestAssign : struct to hold picking assign request data
type rejectRequestAssign struct {
	ID int64 `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	Note               string                    `json:"note"`
	PickingOrderItems  []*approvalPickingItems   `json:"picking_order_items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *rejectRequestAssign) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status != int8(4) {
			o.Failure("picking_order_assign.invalid", util.ErrorPickingSingleStatus("need approval"))
		}
		r.PickingOrderAssign.SalesOrder.Read("ID")
		if r.PickingOrderAssign.SalesOrder.Status != 1 &&
			r.PickingOrderAssign.SalesOrder.Status != 9 &&
			r.PickingOrderAssign.SalesOrder.Status != 12 {
			o.Failure("status.invalid", util.ErrorSalesOrderOnPicking())
		}
		r.PickingOrderAssign.PickingOrder.Read("ID")
		r.PickingOrderAssign.PickingList.Read("ID")
	} else {
		o.Failure("picking_order_assign.invalid", util.ErrorInvalidData("picking order assign"))
	}

	for i, v := range r.PickingOrderItems {
		if pickingOrderAssign, err := common.Decrypt(v.PickingOrderItemID); err == nil {
			if v.PickingOrderItem, err = repository.ValidPickingOrderItem(pickingOrderAssign); err != nil {
				o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
			}
			if v.UnfullfillNote != "" {
				v.POIFlagging = 4
			} else {
				v.POIFlagging = 2
			}

		} else {
			o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *rejectRequestAssign) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
