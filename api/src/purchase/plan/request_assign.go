// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// assignRequest : struct for assignation of Purchase Plan to Field Purchaser
type assignRequest struct {
	ID               int64  `json:"-" valid:"required"`
	FieldPurchaserID string `json:"field_purchaser_id" valid:"required"`

	FieldPurchaser               *model.Staff              `json:"-"`
	PurchasePlan                 *model.PurchasePlan       `json:"-"`
	MessageNotifPreviousAssignee *util.MessageNotification `json:"-"`
	MessageNotifNewAssignee      *util.MessageNotification `json:"-"`
	Session                      *auth.SessionData         `json:"-"`
	PreviousAssignee             *model.Staff              `json:"-"`
}

// Validate : function to validate assign purchase order request data
func (a *assignRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var err error
	var filter, exclude map[string]interface{}

	a.PurchasePlan = &model.PurchasePlan{ID: a.ID}
	if err = a.PurchasePlan.Read("ID"); err != nil {
		o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
	}

	// check if purchase plan status is not active
	if a.PurchasePlan.Status != 1 {
		o.Failure("purchase_plan_status.invalid", util.ErrorActive("purchase plan"))
	}

	filter = map[string]interface{}{"PurchasePlan": a.ID, "Status__in": []int8{1, 2}}
	_, total, err := repository.CheckPurchaseOrderData(filter, exclude)
	if err != nil {
		o.Failure("purchase_order.invalid", util.ErrorInvalidData("purchase order"))
	}

	if total > 0 {
		o.Failure("purchase_plan.invalid", util.ErrorExistActivePurchasePlan())
	}
	// decrypt field purchaser ID
	fieldPurchaserID, err := common.Decrypt(a.FieldPurchaserID)
	if err != nil {
		o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser id"))
	}

	// check if field purchaser ID is valid
	a.FieldPurchaser, err = repository.ValidStaff(fieldPurchaserID)
	if err != nil {
		o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser"))
	}

	if err = a.FieldPurchaser.Role.Read("ID"); err != nil {
		o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser"))
	}

	// check assignee role
	if !(a.FieldPurchaser.Role.Name == "Sourcing Admin" || a.FieldPurchaser.Role.Name == "Field Purchaser") {
		o.Failure("field_purchaser_id.invalid", util.ErrorRole("assignee", "field purchaser or sourcing admin"))
	}

	// check assigner role
	if !(a.Session.Staff.Role.Name == "Sourcing Manager" || a.Session.Staff.Role.Name == "Purchasing Manager") {
		o.Failure("assigner_id.invalid", util.ErrorRole("assigner", "purchasing manager or sourcing manager"))
	}

	if err = a.FieldPurchaser.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	//check if assignee's warehouse is same with assigner's warehouse
	if a.FieldPurchaser.Warehouse.Name != "All Warehouse" {
		if a.FieldPurchaser.Warehouse.ID != a.PurchasePlan.Warehouse.ID {
			o.Failure("warehouse.invalid", util.ErrorMustBeSame("Field Purchaser Warehouse", "Purchase Plan Warehouse"))
		}
	}

	if err = a.PurchasePlan.SupplierOrganization.Read("ID"); err != nil {
		o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization"))
	}

	if err = a.FieldPurchaser.User.Read("ID"); err != nil {
		o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
	}

	if a.PurchasePlan.AssignedTo != nil {
		a.PreviousAssignee, err = repository.ValidStaff(a.PurchasePlan.AssignedTo.ID)
		if err != nil {
			o.Failure("staff_id.invalid", util.ErrorInvalidData("staff id"))
		}

		if err = a.PreviousAssignee.User.Read("ID"); err != nil {
			o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
		}

		if a.FieldPurchaser.ID != a.PreviousAssignee.ID {
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0024'").QueryRow(&a.MessageNotifPreviousAssignee)
			a.MessageNotifPreviousAssignee.Message = util.ReplaceCodeString(a.MessageNotifPreviousAssignee.Message, map[string]interface{}{"#purchase_plan_code#": a.PurchasePlan.Code, "#purchasing_manager_name#": a.Session.Staff.DisplayName})

			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0023'").QueryRow(&a.MessageNotifNewAssignee)
			a.MessageNotifNewAssignee.Message = util.ReplaceCodeString(a.MessageNotifNewAssignee.Message, map[string]interface{}{"#purchase_plan_code#": a.PurchasePlan.Code})
		}
	} else {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0023'").QueryRow(&a.MessageNotifNewAssignee)
		a.MessageNotifNewAssignee.Message = util.ReplaceCodeString(a.MessageNotifNewAssignee.Message, map[string]interface{}{"#purchase_plan_code#": a.PurchasePlan.Code})
	}

	return o
}

// Messages : function to return error validation messages
func (a *assignRequest) Messages() map[string]string {
	return map[string]string{
		"id.required":                 util.ErrorInputRequired("id"),
		"field_purchaser_id.required": util.ErrorInputRequired("field purchaser"),
	}
}
