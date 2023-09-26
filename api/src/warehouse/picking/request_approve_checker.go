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

// approveRequestChecker : struct to hold picking assign request data
type approveRequestChecker struct {
	ID int64 `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	TotalColly         float64                   `json:"total_colly"`
	PickingOrderItems  []*approvalPickingItems   `json:"picking_order_items"`
	DeliveryKolies     []*DeliveryKoli           `json:"delivery_kolies" valid:"required"`

	IsFinished bool `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *approveRequestChecker) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.TotalColly <= 0 {
		o.Failure("id.invalid", util.ErrorEqualGreaterInd("total koli", "0"))
		return o
	}

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status != int8(6) {
			o.Failure("picking_order_assign.invalid", util.ErrorPickingSingleStatus("checking"))
		}
		r.PickingOrderAssign.SalesOrder.Read("ID")
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
func (r *approveRequestChecker) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
