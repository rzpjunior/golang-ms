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

type cancelRequest struct {
	ID   int64  `json:"-" valid:"required"`
	Note string `json:"note" valid:"required|lte:250"`

	PurchasePlan          *model.PurchasePlan       `json:"-"`
	Session               *auth.SessionData         `json:"-"`
	MessageNotifPurchaser *util.MessageNotification `json:"-"`
	MessageNotifManager   *util.MessageNotification `json:"-"`
	FieldPurchaser        *model.Staff              `json:"-"`
	PurchasingManagers    []*model.Staff            `json:"-"`
}

func (r *cancelRequest) Validate() *validation.Output {
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

	filter = map[string]interface{}{"Role__Name__in": []string{"Purchasing Manager", "Sourcing Manager"}, "Warehouse__ID__in": []int64{r.PurchasePlan.Warehouse.ID, 21}, "Status": 1}
	r.PurchasingManagers, _, err = repository.CheckStaffData(filter, exclude)
	if err != nil {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	if r.PurchasePlan.AssignedTo != nil {

		r.FieldPurchaser, err = repository.ValidStaff(r.PurchasePlan.AssignedTo.ID)
		if err != nil {
			o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field_purchaser"))
		}

		if err = r.FieldPurchaser.User.Read("ID"); err != nil {
			o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
		}

		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0026'").QueryRow(&r.MessageNotifPurchaser)
		r.MessageNotifPurchaser.Message = util.ReplaceCodeString(r.MessageNotifPurchaser.Message, map[string]interface{}{"#purchase_plan_code#": r.PurchasePlan.Code})
	}

	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0026'").QueryRow(&r.MessageNotifManager)
	r.MessageNotifManager.Message = util.ReplaceCodeString(r.MessageNotifManager.Message, map[string]interface{}{"#purchase_plan_code#": r.PurchasePlan.Code})

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
		"note.lte":      util.ErrorEqualLess("cancellation note", "250"),
	}
}
