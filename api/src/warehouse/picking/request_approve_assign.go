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

// approveRequestAssign : struct to hold picking assign request data
type approveRequestAssign struct {
	ID int64 `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	PickingOrderItems  []*approvalPickingItems   `json:"picking_order_items" valid:"required"`
	IsFinished         bool                      `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type DeliveryKoli struct {
	KoliID   string  `json:"koli_id"`
	Quantity float64 `json:"quantity"`
	Note     string  `json:"note"`

	SalesOrder *model.SalesOrder `json:"-"`
	Koli       *model.Koli       `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *approveRequestAssign) Validate() *validation.Output {
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
	} else {
		o.Failure("picking_order_assign.invalid", util.ErrorInvalidData("picking order assign"))
	}

	for i, v := range r.PickingOrderItems {
		if pickingOrderAssign, err := common.Decrypt(v.PickingOrderItemID); err == nil {
			if v.PickingOrderItem, err = repository.ValidPickingOrderItem(pickingOrderAssign); err != nil {
				o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
			}
		} else {
			o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *approveRequestAssign) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
