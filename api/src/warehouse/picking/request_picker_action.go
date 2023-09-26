// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// requestPickerAction : struct to hold picking assign request data
type requestPickerAction struct {
	ID            int64   `json:"-"`
	PickOrderQty  float64 `json:"pick_qty"`
	UnfulfillNote string  `json:"unfulfill_note"`
	PickingFlag   int8    `json:"-"`

	PickingList        *model.PickingList        `json:"-"`
	PickingRoutingStep *model.PickingRoutingStep `json:"-"`
	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	PickingOrderItem   *model.PickingOrderItem   `json:"-"`

	ActionType             *model.Glossary             `json:"-"`
	IsPicking              bool                        `json:"-"`
	NextPickingRoutingStep *model.PickingRoutingStep   `json:"-"`
	AllPickingRoutingStep  []*model.PickingRoutingStep `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *requestPickerAction) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var filter, exclude map[string]interface{}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	//check the picking list
	r.PickingList = &model.PickingList{ID: r.ID}
	err = r.PickingList.Read("id")
	if err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}

	// get the picking routing step
	if err = o1.Raw("SELECT * FROM picking_routing_step WHERE staff_id = ? and picking_list_id = ? and status_step = 2 ORDER BY sequence ASC LIMIT 1", r.Session.Staff.ID, r.ID).QueryRow(&r.PickingRoutingStep); err != nil {
		o.Failure("picking_routing_step.invalid", util.ErrorInvalidData("picking routing step"))
		return o
	}

	// determine the type from glossary
	r.ActionType, err = repository.GetGlossaryMultipleValue("table", "picking_routing_step", "attribute", "step_type", "value_int", r.PickingRoutingStep.StepType)
	if err != nil {
		o.Failure("step_type_id.invalid", util.ErrorInvalidData("picking routing step"))
		return o
	}

	// if the staff already walked to the location, then he's trying to pick
	if r.PickingRoutingStep.WalkingFinishTime.IsZero() && r.ActionType.ValueName == "pickup" {
	} else {
		r.IsPicking = true
	}

	if r.ActionType.ValueName == "pickup" && r.IsPicking == true {
		r.PickingOrderItem = &model.PickingOrderItem{ID: r.PickingRoutingStep.PickingOrderItem.ID}
		if err = r.PickingOrderItem.Read("id"); err != nil {
			o.Failure("picking_order_item_id.invalid", util.ErrorInvalidData("picking order item"))
			return o
		}

		if err = r.PickingOrderItem.PickingOrderAssign.Read("id"); err != nil {
			o.Failure("picking_order_assign_id.invalid", util.ErrorInvalidData("picking order assign"))
		}
		if err = r.PickingOrderItem.PickingOrderAssign.SalesOrder.Read("id"); err != nil {
			o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		}

		if r.PickingOrderItem.PickingOrderAssign.SalesOrder.Status != 3 {
			if r.PickOrderQty < r.PickingOrderItem.OrderQuantity && r.UnfulfillNote == "" {
				o.Failure("Unfulfill_note.invalid", util.ErrorInputRequired("unfulfill note"))
				return o
			}
			if r.UnfulfillNote != "" {
				r.PickingFlag = 3
			} else {
				r.PickingFlag = 2
			}

		}
	}

	if r.ActionType.ValueName == "end" {
		filter = map[string]interface{}{"staff_id": r.PickingRoutingStep.Staff.ID, "status_step__in": []int{2, 3}}
		r.AllPickingRoutingStep, _, err = repository.CheckPickingRoutingStepData(filter, exclude)
		if err != nil {
			o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		}
	} else {
		if err = o1.Raw("SELECT * FROM picking_routing_step WHERE staff_id = ? and picking_list_id = ? and sequence = ? and status_step=2", r.PickingRoutingStep.Staff.ID, r.PickingRoutingStep.PickingList.ID, (r.PickingRoutingStep.Sequence + 1)).QueryRow(&r.NextPickingRoutingStep); err != nil {
			o.Failure("picking_routing_step.invalid", util.ErrorInvalidData("picking routing step"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *requestPickerAction) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
