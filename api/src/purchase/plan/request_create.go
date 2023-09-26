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

// createRequest : struct to hold Create Purchase Plan request data

type createRequest struct {
	Code                   string    `json:"-"`
	SupplierOrganizationID string    `json:"supplier_organization_id" valid:"required"`
	WarehouseID            string    `json:"warehouse_id" valid:"required"`
	StrRecognitionDate     string    `json:"recognition_date" valid:"required"`
	StrEtaDate             string    `json:"eta_date" valid:"required"`
	EtaTime                string    `json:"eta_time" valid:"required"`
	FieldPurchaserID       string    `json:"field_purchaser_id"`
	Note                   string    `json:"note" valid:"lte:250"`
	RecognitionDate        time.Time `json:"-"`
	EtaDate                time.Time `json:"-"`

	TotalPurchasePlanQty float64   `json:"-"`
	EtaTimeFormat        time.Time `json:"-"`

	PurchasePlanItems []*requestItem `json:"purchase_plan_items" valid:"required"`

	TotalPrice            float64                     `json:"-"`
	TotalWeight           float64                     `json:"-"`
	RecognitionAt         time.Time                   `json:"-"`
	EtaDateAt             time.Time                   `json:"-"`
	SupplierOrganization  *model.SupplierOrganization `json:"-"`
	Warehouse             *model.Warehouse            `json:"-"`
	FieldPurchaser        *model.Staff                `json:"-"`
	MessageNotifPurchaser *util.MessageNotification   `json:"-"`
	MessageNotifManager   *util.MessageNotification   `json:"-"`
	PurchasingManagers    []*model.Staff              `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type requestItem struct {
	ID              string  `json:"id"`
	ProductID       string  `json:"product_id" valid:"required"`
	PurchasePlanQty float64 `json:"purchase_plan_qty" valid:"required|range:0,99999999.99"`
	UnitPrice       float64 `json:"unit_price" valid:"required|range:0,9999999999.99"`

	TaxableItem      int8                    `json:"-"`
	TaxAmount        float64                 `json:"-"`
	UnitPriceTax     float64                 `json:"-"`
	Subtotal         float64                 `json:"-"`
	PurchasePlanItem *model.PurchasePlanItem `json:"-"`
	Product          *model.Product          `json:"-"`
	Price            *model.Price            `json:"-"`
	Uom              *model.Uom              `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var filter, exclude map[string]interface{}
	productList := make(map[int64]string)

	supplierOrganizationID, err := common.Decrypt(c.SupplierOrganizationID)
	if err != nil {
		o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization id"))
	}

	if c.Code, err = util.CheckTable("purchase_plan"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	c.SupplierOrganization, err = repository.ValidSupplierOrganization(supplierOrganizationID)
	if err != nil {
		o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization id"))
	}

	warehouseID, e := common.Decrypt(c.WarehouseID)
	if e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	c.Warehouse, err = repository.ValidWarehouse(warehouseID)
	if err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if c.FieldPurchaserID != "" {
		fieldPurchaserID, e := common.Decrypt(c.FieldPurchaserID)
		if e != nil {
			o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field purchaser"))
		}

		c.FieldPurchaser, err = repository.ValidStaff(fieldPurchaserID)
		if err != nil {
			o.Failure("field_purchaser_id.invalid", util.ErrorInvalidData("field_purchaser"))
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
			if c.FieldPurchaser.Warehouse.ID != c.Warehouse.ID {
				o.Failure("warehouse.invalid", util.ErrorMustBeSame("Field Purchaser Warehouse", "Purchase Plan Warehouse"))
			}
		}

		if err = c.FieldPurchaser.User.Read("ID"); err != nil {
			o.Failure("user_id.invalid", util.ErrorInvalidData("user id"))
		}
	}

	filter = map[string]interface{}{"Role__Name__in": []string{"Purchasing Manager", "Sourcing Manager"}, "Warehouse__ID__in": []int64{c.Warehouse.ID, 21}, "Status": 1}
	c.PurchasingManagers, _, err = repository.CheckStaffData(filter, exclude)
	if err != nil {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	if c.FieldPurchaserID != "" {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0023'").QueryRow(&c.MessageNotifPurchaser)
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0027'").QueryRow(&c.MessageNotifManager)
	} else {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0022'").QueryRow(&c.MessageNotifManager)
	}

	if c.RecognitionDate, err = time.Parse("2006-01-02", c.StrRecognitionDate); err != nil {
		o.Failure("order_date.invalid", util.ErrorInvalidData("order date"))
	}

	if c.EtaDate, err = time.Parse("2006-01-02", c.StrEtaDate); err != nil {
		o.Failure("eta_date.invalid", util.ErrorInvalidData("estimated arrival date"))
	}

	// only for checking format time from apps
	if c.EtaTimeFormat, e = time.Parse("15:04", c.EtaTime); e != nil {
		o.Failure("eta_time.invalid", util.ErrorInvalidData("eta time"))
	}

	for i, v := range c.PurchasePlanItems {

		if v.PurchasePlanQty <= 0 {
			o.Failure("order_qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("product quantity", "0"))
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

		v.Subtotal = v.PurchasePlanQty * v.UnitPrice
		c.TotalPurchasePlanQty += v.PurchasePlanQty
		c.TotalPrice = c.TotalPrice + v.Subtotal
		c.TotalWeight = c.TotalWeight + (v.PurchasePlanQty * v.Product.UnitWeight)

		if _, exist := productList[productID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
		}

		productList[productID] = "y"

		filter = map[string]interface{}{"product_id": productID, "warehouse_id": warehouseID, "purchasable": 1}
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"supplier_organization_id.required": util.ErrorSelectRequired("supplier organization"),
		"warehouse_id.required":             util.ErrorSelectRequired("warehouse"),
		"recognition_date.required":         util.ErrorInputRequired("purchase plan date"),
		"eta_date.required":                 util.ErrorInputRequired("eta date"),
		"eta_time.required":                 util.ErrorInputRequired("eta time"),
		"note.lte":                          util.ErrorEqualLess("note", "250"),
		"purchase_plan_items.required":      util.ErrorInputRequired("product item"),
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
