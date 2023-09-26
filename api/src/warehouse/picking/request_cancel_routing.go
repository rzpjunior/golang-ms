// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// startRoutingAssignment : struct to hold picking assign routing request data
type cancelRoutingAssignment struct {
	ID int64 `json:"-"`

	PickingList        *model.PickingList          `json:"-"`
	PickingRoutingStep []*model.PickingRoutingStep `json:"-"`
	PickingOrderAssign []*model.PickingOrderAssign `json:"-"`
	Picker             []int64                     `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign routing request data
func (r *cancelRoutingAssignment) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err             error
		filter, exclude map[string]interface{}
		total           int64
	)
	pickerMap := map[int64]bool{}

	//check the picking list
	r.PickingList = &model.PickingList{ID: r.ID}
	err = r.PickingList.Read("id")
	if err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}

	filter = map[string]interface{}{"picking_list_id": r.ID}
	r.PickingOrderAssign, _, err = repository.CheckPickingOrderAssignData(filter, exclude)
	if err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
	}

	for _, v := range r.PickingOrderAssign {
		if v.Helper.ID != r.Session.Staff.ID {
			o.Failure("staff_id.invalid", util.ErrorPickingListStaff())
		}
	}

	// get all picking routing step
	filter = map[string]interface{}{"picking_list_id": r.ID, "status_step__in": []int64{2, 3}}
	r.PickingRoutingStep, total, err = repository.CheckPickingRoutingStepData(filter, exclude)
	if err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
	}
	if total == 0 {
		o.Failure("picking_list_id.invalid", util.ErrorNotFound("picking routing step"))
	}

	for _, v := range r.PickingRoutingStep {
		if pickerMap[v.Staff.ID] == false {
			pickerMap[v.Staff.ID] = true
			r.Picker = append(r.Picker, v.Staff.ID)
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *cancelRoutingAssignment) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
