// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type cancelAssignmentRequest struct {
	ID int64 `json:"-" valid:"required"`

	PurchasePlan          *model.PurchasePlan       `json:"-"`
	Session               *auth.SessionData         `json:"-"`
	MessageNotifPurchaser *util.MessageNotification `json:"-"`
	MessageNotifManager   *util.MessageNotification `json:"-"`
	PreviousAssignee      *model.Staff              `json:"-"`
}

func (r *cancelAssignmentRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var err error
	var filter, exclude map[string]interface{}

	r.PurchasePlan, err = repository.ValidPurchasePlan(r.ID)
	if err != nil {
		o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
	}

	if r.PurchasePlan.Status != 1 {
		o.Failure("purchase_plan.invalid", util.ErrorActive("purchase plan"))
	}

	filter = map[string]interface{}{"PurchasePlan": r.ID, "Status__in": []int8{1, 2}}
	_, total, err := repository.CheckPurchaseOrderData(filter, exclude)
	if err != nil {
		o.Failure("purchase_order.invalid", util.ErrorInvalidData("purchase order"))
	}

	if total > 0 {
		o.Failure("purchase_plan.invalid", util.ErrorExistActivePurchasePlan())
	}

	r.PreviousAssignee, err = repository.ValidStaff(r.PurchasePlan.AssignedTo.ID)
	if err != nil {
		o.Failure("staff_id.invalid", util.ErrorInvalidData("staff"))
	}

	if err = r.PreviousAssignee.User.Read("ID"); err != nil {
		o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
	}

	if err = r.Session.Staff.User.Read("ID"); err != nil {
		o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
	}

	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0024'").QueryRow(&r.MessageNotifPurchaser)
	r.MessageNotifPurchaser.Message = util.ReplaceCodeString(r.MessageNotifPurchaser.Message, map[string]interface{}{"#purchase_plan_code#": r.PurchasePlan.Code, "#purchasing_manager_name#": r.Session.Staff.DisplayName})

	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0025'").QueryRow(&r.MessageNotifManager)
	r.MessageNotifManager.Message = util.ReplaceCodeString(r.MessageNotifManager.Message, map[string]interface{}{"#purchase_plan_code#": r.PurchasePlan.Code, "#field_purchaser_name#": r.PreviousAssignee.DisplayName})

	return o
}

func (r *cancelAssignmentRequest) Messages() map[string]string {
	return map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}
}
