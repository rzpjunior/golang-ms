package order

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// assignRequest : struct for assignation of Purchase Order to Field Purchaser
type assignRequest struct {
	ID               int64  `json:"-" valid:"required"`
	FieldPurchaserID string `json:"field_purchaser_id" valid:"required"`

	FieldPurchaser *model.Staff              `json:"-"`
	PurchaseOrder  *model.PurchaseOrder      `json:"-"`
	MessageNotif   *util.MessageNotification `json:"-"`
	Session        *auth.SessionData         `json:"-"`
}

// Validate : function to validate assign purchase order request data
func (a *assignRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var err error

	a.PurchaseOrder = &model.PurchaseOrder{ID: a.ID}
	if err = a.PurchaseOrder.Read("ID"); err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	// check if purchase order status is not draft
	if a.PurchaseOrder.Status != 5 {
		o.Failure("status.invalid", util.ErrorDraft("purchase order"))
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

	if err = a.FieldPurchaser.Read("ID"); err != nil {
		o.Failure("staff_id.invalid", util.ErrorInvalidData("staff id"))
	}

	if err = a.FieldPurchaser.Role.Read("ID"); err != nil {
		o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser"))
	}

	// check assignee role
	if !(a.FieldPurchaser.Role.Name == "Sourcing Admin" || a.FieldPurchaser.Role.Name == "Field Purchaser") {
		o.Failure("field_purchaser_id.invalid", util.ErrorRole("field purchaser", "field purchaser or sourcing admin"))
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
		if a.FieldPurchaser.Warehouse.ID != a.PurchaseOrder.Warehouse.ID {
			o.Failure("warehouse.invalid", util.ErrorMustBeSame("Field Purchaser Warehouse", "Purchase Order Warehouse"))
		}
	}

	if err = a.PurchaseOrder.Supplier.Read("ID"); err != nil {
		o.Failure("supplier_id.invalid", util.ErrorInvalidData("supplier"))
	}

	if err = a.FieldPurchaser.User.Read("ID"); err != nil {
		o.Failure("user_id.invalid", util.ErrorInvalidData("user_id"))
	}

	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0019'").QueryRow(&a.MessageNotif)
	a.MessageNotif.Title = util.ReplaceCodeString(a.MessageNotif.Title, map[string]interface{}{"#supplier_name#": a.PurchaseOrder.Supplier.Name})
	a.MessageNotif.Message = util.ReplaceCodeString(a.MessageNotif.Message, map[string]interface{}{"#name#": a.FieldPurchaser.DisplayName, "#smile#": "☺"})
	return o
}

// Messages : function to return error validation messages
func (a *assignRequest) Messages() map[string]string {
	return map[string]string{
		"id.required":                 util.ErrorInputRequired("id"),
		"field_purchaser_id.required": util.ErrorInputRequired("field purchaser"),
	}
}
