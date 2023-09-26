// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consolidated_shipment

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// consolidateRequest : struct to hold input consolidate shipment request data
type consolidateRequest struct {
	Code              string           `json:"-"`
	DriverName        string           `json:"driver_name" valid:"required|alpha_num_space|lte:100"`
	VehicleNumber     string           `json:"vehicle_number" valid:"required|alpha_num|range:5,9"`
	DriverPhoneNumber string           `json:"driver_phone_number" valid:"required|numeric|range:8,15"`
	PurchaseOrders    []*purchaseOrder `json:"purchase_orders" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type purchaseOrder struct {
	PurchaseOrderID string `json:"purchase_order_id" valid:"required"`

	PurchaseOrder *model.PurchaseOrder `json:"-"`
}

// Validate : function to validate create field purchase order request data
func (c *consolidateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var warehouse, supplierOrganization []string

	if c.Code, err = util.CheckTable("consolidated_shipment"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	for k, v := range c.PurchaseOrders {
		purchaseOrderID, err := common.Decrypt(v.PurchaseOrderID)
		if err != nil {
			o.Failure("purchase_order_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("purchase order id"))
		}

		v.PurchaseOrder = &model.PurchaseOrder{ID: purchaseOrderID}
		if err = v.PurchaseOrder.Read("ID"); err != nil {
			o.Failure("purchase_order_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("purchase order id"))
		}

		if err = v.PurchaseOrder.PurchasePlan.Read("ID"); err != nil {
			o.Failure("purchase_plan_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("purchase plan id"))
		}

		if v.PurchaseOrder.ConsolidatedShipment != nil {
			o.Failure("purchase_order"+strconv.Itoa(k)+".invalid", util.ErrorSelectAnother("purchase order", "consolidated", "purchase order"))
		}

		warehouse = append(warehouse, strconv.FormatInt(v.PurchaseOrder.Warehouse.ID, 10))
		supplierOrganization = append(supplierOrganization, strconv.FormatInt(v.PurchaseOrder.PurchasePlan.SupplierOrganization.ID, 10))
	}

	warehouse = util.RemoveDuplicateValuesString(warehouse)
	supplierOrganization = util.RemoveDuplicateValuesString(supplierOrganization)

	if len(warehouse) > 1 {
		o.Failure("warehouse.invalid", util.ErrorMustBeSameInOneDocument("warehouse", "consolidated shipment"))
	}

	if len(supplierOrganization) > 1 {
		o.Failure("supplier_organization.invalid", util.ErrorMustBeSameInOneDocument("supplier organization", "consolidated shipment"))
	}

	// get warehouse code
	if err := c.PurchaseOrders[0].PurchaseOrder.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse id"))
	}

	c.DriverPhoneNumber = util.ParsePhoneNumberPrefix(c.DriverPhoneNumber)

	return o
}

func (c *consolidateRequest) Messages() map[string]string {
	messages := map[string]string{
		"driver_name.required":         util.ErrorInputRequired("driver name"),
		"vehicle_number.required":      util.ErrorInputRequired("vehicle number"),
		"purchase_orders.required":     util.ErrorSelectRequired("purchase order"),
		"purchase_order_id.required":   util.ErrorInputRequired("purchase order id"),
		"driver_name.alpha_num_space":  util.ErrorAlphaNum("driver name"),
		"driver_name.lte":              util.ErrorEqualLess("driver name", "100"),
		"vehicle_number.alpha_num":     util.ErrorAlphaNum("vehicle number"),
		"vehicle_number.range":         util.ErrorRangeChar("vehicle number", "5", "9"),
		"driver_phone_number.numeric":  util.ErrorNumeric("driver phone number"),
		"driver_phone_number.range":    util.ErrorRangeChar("driver phone number", "8", "15"),
		"driver_phone_number.required": util.ErrorInputRequired("driver phone number"),
	}

	for i, _ := range c.PurchaseOrders {
		messages["item."+strconv.Itoa(i)+".purchase_order_id.required"] = util.ErrorSelectRequired("purchase order")
	}

	return messages
}
