// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"math"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequest : struct to hold Update Purchase Plan request data
type updateRequest struct {
	ID                 int64     `json:"-"`
	StrRecognitionDate string    `json:"recognition_date" valid:"required"`
	StrEtaDate         string    `json:"eta_date" valid:"required"`
	EtaTime            string    `json:"eta_time" valid:"required"`
	FieldPurchaserID   string    `json:"field_purchaser_id"`
	Note               string    `json:"note" valid:"lte:250"`
	RecognitionDate    time.Time `json:"-"`
	EtaDate            time.Time `json:"-"`

	TotalPurchasePlanQty float64   `json:"-"`
	EtaTimeFormat        time.Time `json:"-"`

	PurchasePlanItems []*requestItem `json:"purchase_plan_items" valid:"required"`

	TotalPrice                   float64                   `json:"-"`
	TotalWeight                  float64                   `json:"-"`
	RecognitionAt                time.Time                 `json:"-"`
	EtaDateAt                    time.Time                 `json:"-"`
	PurchasePlan                 *model.PurchasePlan       `json:"-"`
	FieldPurchaser               *model.Staff              `json:"-"`
	MessageNotifPreviousAssignee *util.MessageNotification `json:"-"`
	MessageNotifNewAssignee      *util.MessageNotification `json:"-"`
	MessageNotifManager          *util.MessageNotification `json:"-"`
	Session                      *auth.SessionData         `json:"-"`
	PurchasingManagers           []*model.Staff            `json:"-"`
	PreviousAssignee             *model.Staff              `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var err error
	var filter, exclude map[string]interface{}
	productList := make(map[int64]string)

	c.PurchasePlan, err = repository.ValidPurchasePlan(c.ID)
	if err != nil {
		o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
	}

	if c.PurchasePlan.Status != 1 {
		o.Failure("purchase_plan.invalid", util.ErrorActive("purchase plan"))
	}

	if c.RecognitionDate, err = time.Parse("2006-01-02", c.StrRecognitionDate); err != nil {
		o.Failure("purchase_plan_date.invalid", util.ErrorInvalidData("purchase plan date"))
	}

	if c.EtaDate, err = time.Parse("2006-01-02", c.StrEtaDate); err != nil {
		o.Failure("eta_date.invalid", util.ErrorInvalidData("estimated arrival date"))
	}

	// only for checking format time from apps
	if c.EtaTimeFormat, err = time.Parse("15:04", c.EtaTime); err != nil {
		o.Failure("eta_time.invalid", util.ErrorInvalidData("eta time"))
	}

	if c.FieldPurchaserID != "" {
		fieldPurchaserID, e := common.Decrypt(c.FieldPurchaserID)
		if e != nil {
			o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser"))
		}

		c.FieldPurchaser, err = repository.ValidStaff(fieldPurchaserID)
		if err != nil {
			o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser"))
		}

		if err = c.FieldPurchaser.Role.Read("ID"); err != nil {
			o.Failure("field_purchaser_role_id.invalid", util.ErrorInvalidData("field purchaser role id"))
		}

		// check assignee role
		if !(c.FieldPurchaser.Role.Name == "Field Purchaser" || c.FieldPurchaser.Role.Name == "Sourcing Admin") {
			o.Failure("field_purchaser_id.invalid", util.ErrorRole("assignee", "field purchaser or sourcing admin"))
		}

		if err = c.FieldPurchaser.Warehouse.Read("ID"); err != nil {
			o.Failure("field_purchaser_warehouse_id.invalid", util.ErrorInvalidData("field purchaser warehouse id"))
		}

		//check if assignee's warehouse is same with assigner's warehouse
		if c.FieldPurchaser.Warehouse.Name != "All Warehouse" {
			if c.FieldPurchaser.Warehouse.ID != c.PurchasePlan.Warehouse.ID {
				o.Failure("warehouse.invalid", util.ErrorMustBeSame("Field Purchaser Warehouse", "Purchase Plan Warehouse"))
			}
		}

		if err = c.PurchasePlan.SupplierOrganization.Read("ID"); err != nil {
			o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization"))
		}

		if err = c.FieldPurchaser.User.Read("ID"); err != nil {
			o.Failure("user_id.invalid", util.ErrorInvalidData("user_id"))
		}

		filter = map[string]interface{}{"PurchasePlan": c.ID, "Status": 1}
		_, total, err := repository.CheckPurchaseOrderData(filter, exclude)
		if err != nil {
			o.Failure("purchase_order.invalid", util.ErrorInvalidData("purchase order"))
		}

		if total > 0 {
			if c.FieldPurchaser.ID != c.PurchasePlan.AssignedTo.ID {
				o.Failure("field_purchaser.invalid", util.ErrorCannotUpdateAfter("Field Purchaser", "Create Purchase Order"))
			}
		}
	}

	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0023'").QueryRow(&c.MessageNotifNewAssignee)
	c.MessageNotifNewAssignee.Message = util.ReplaceCodeString(c.MessageNotifNewAssignee.Message, map[string]interface{}{"#purchase_plan_code#": c.PurchasePlan.Code})

	if c.PurchasePlan.AssignedTo != nil {
		if c.FieldPurchaserID == "" {
			filter = map[string]interface{}{"Role__Name__in": []string{"Purchasing Manager", "Sourcing Manager"}, "Warehouse__ID__in": []int64{c.PurchasePlan.Warehouse.ID, 21}, "Status": 1}
			c.PurchasingManagers, _, err = repository.CheckStaffData(filter, exclude)
			if err != nil {
				o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
			}
		}

		c.PreviousAssignee, err = repository.ValidStaff(c.PurchasePlan.AssignedTo.ID)
		if err != nil {
			o.Failure("staff_id.invalid", util.ErrorInvalidData("staff id"))
		}

		if err = c.PreviousAssignee.User.Read("ID"); err != nil {
			o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
		}

		if c.FieldPurchaserID == "" {
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0024'").QueryRow(&c.MessageNotifPreviousAssignee)
			c.MessageNotifPreviousAssignee.Message = util.ReplaceCodeString(c.MessageNotifPreviousAssignee.Message, map[string]interface{}{"#purchase_plan_code#": c.PurchasePlan.Code, "#purchasing_manager_name#": c.Session.Staff.DisplayName})

			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0025'").QueryRow(&c.MessageNotifManager)
			c.MessageNotifManager.Message = util.ReplaceCodeString(c.MessageNotifManager.Message, map[string]interface{}{"#purchase_plan_code#": c.PurchasePlan.Code, "#field_purchaser_name#": c.PreviousAssignee.DisplayName})
		} else if c.FieldPurchaser.ID != c.PreviousAssignee.ID {
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0024'").QueryRow(&c.MessageNotifPreviousAssignee)
			c.MessageNotifPreviousAssignee.Message = util.ReplaceCodeString(c.MessageNotifPreviousAssignee.Message, map[string]interface{}{"#purchase_plan_code#": c.PurchasePlan.Code, "#purchasing_manager_name#": c.Session.Staff.DisplayName})

			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0023'").QueryRow(&c.MessageNotifNewAssignee)
			c.MessageNotifNewAssignee.Message = util.ReplaceCodeString(c.MessageNotifNewAssignee.Message, map[string]interface{}{"#purchase_plan_code#": c.PurchasePlan.Code})
		}
	}

	for i, v := range c.PurchasePlanItems {

		if v.PurchasePlanQty <= 0 {
			o.Failure("qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("product quantity", "0"))
		}

		productID, err := common.Decrypt(v.ProductID)
		if err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}

		v.Product, err = repository.ValidProduct(productID)
		if err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}

		v.Uom, err = repository.ValidUom(v.Product.Uom.ID)
		if err != nil {
			o.Failure("uom_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("uom"))
		}

		if v.Uom.DecimalEnabled == 2 {
			if math.Mod(v.PurchasePlanQty, 1) != 0 {
				o.Failure("order_qty"+strconv.Itoa(i)+".invalid", util.ErrorNotAllowedFor("decimal", "product qty"))
			}
		}

		if v.ID != "" {
			purchasePlanItemID, err := common.Decrypt(v.ID)
			if err != nil {
				o.Failure("purchase_plan_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase plan item"))
			}

			if v.PurchasePlanItem, err = repository.ValidPurchasePlanItem(purchasePlanItemID); err != nil {
				o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			}

			purchasePlanItemID, e := common.Decrypt(v.ID)
			if e != nil {
				o.Failure("purchase_plan_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase plan item"))
			}

			v.PurchasePlanItem = &model.PurchasePlanItem{ID: purchasePlanItemID}
			err = v.PurchasePlanItem.Read("ID")
			if err != nil {
				o.Failure("purchase_plan_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase plan item"))
			}

			if v.PurchasePlanItem.PurchaseQty > v.PurchasePlanQty {
				o.Failure("purchase_plan_qty"+strconv.Itoa(i)+".invalid", util.ErrorEqualGreater("purchase plan quantity", "purchase quantity"))
			}

			if v.PurchasePlanItem.PurchaseQty > 0 {
				if v.PurchasePlanItem.Product.ID != v.Product.ID {
					o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorCannotUpdateAfter("product", "purchase order created for the selected purchase plan"))
				}
			}
		}

		v.Subtotal = v.PurchasePlanQty * v.UnitPrice
		c.TotalPurchasePlanQty += v.PurchasePlanQty
		c.TotalPrice = c.TotalPrice + v.Subtotal
		c.TotalWeight = c.TotalWeight + (v.PurchasePlanQty * v.Product.UnitWeight)

		if _, exist := productList[productID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
		}

		productList[productID] = "y"

		filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.PurchasePlan.Warehouse.ID, "purchasable": 1}
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}
	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"recognition_date.required":    util.ErrorInputRequired("purchase plan date"),
		"eta_date.required":            util.ErrorInputRequired("eta date"),
		"eta_time.required":            util.ErrorInputRequired("eta time"),
		"note.lte":                     util.ErrorEqualLess("note", "250"),
		"purchase_plan_items.required": util.ErrorInputRequired("product item"),
	}

	for i, _ := range c.PurchasePlanItems {
		messages["item."+strconv.Itoa(i)+".product_id.required"] = util.ErrorSelectRequired("product")
		messages["item."+strconv.Itoa(i)+".purchase_plan_qty.required"] = util.ErrorInputRequired("purchase plan qty")
		messages["item."+strconv.Itoa(i)+".unit_price.required"] = util.ErrorInputRequired("unit price")
		messages["item."+strconv.Itoa(i)+".purchase_plan_qty.range"] = util.ErrorRangeValue("purchase plan qty", "0", "99999999.99")
		messages["item."+strconv.Itoa(i)+".unit_price.range"] = util.ErrorRangeValue("unit price", "0", "9999999999,99")
	}

	return messages
}
