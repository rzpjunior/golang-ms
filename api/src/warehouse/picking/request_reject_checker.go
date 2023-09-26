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

// rejectRequestChecker : struct to hold picking assign request data
type rejectRequestChecker struct {
	ID int64 `json:"-"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	Note               string                    `json:"note"`
	PickingOrderItems  []*approvalPickingItems   `json:"picking_order_items"`

	CanceledPickingOrderItemsID  []int64                     `json:"-"`
	DeletedPOIPickingRoutingStep []*model.PickingRoutingStep `json:"-"`

	IsFinished bool `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *rejectRequestChecker) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var filter, exclude map[string]interface{}
	var err error

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err == nil {
		if r.PickingOrderAssign.Status != int8(6) {
			o.Failure("picking_order_assign.invalid", util.ErrorPickingSingleStatus("checking"))
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

			if _, err := repository.GetGlossaryMultipleValue("table", "picking_order_item", "attribute", "picking_flag", "value_int", v.POIFlagging); err != nil {
				o.Failure("glossary_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("glossary"))
			}
		} else {
			o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
		}
	}

	flagOrder := 4
	o1.Raw("Select id from picking_order_item where picking_order_assign_id = ? and flag_order = ?", r.PickingOrderAssign.ID, flagOrder).QueryRows(&r.CanceledPickingOrderItemsID)
	if len(r.CanceledPickingOrderItemsID) > 0 {
		filter = map[string]interface{}{"picking_order_item_id__in": r.CanceledPickingOrderItemsID}
		if r.DeletedPOIPickingRoutingStep, _, err = repository.CheckPickingRoutingStepData(filter, exclude); err != nil {
			o.Failure("picking_routing_step.invalid", util.ErrorInvalidData("picking routing step"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *rejectRequestChecker) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
