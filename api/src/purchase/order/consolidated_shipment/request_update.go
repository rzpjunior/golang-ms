// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consolidated_shipment

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

// updateRequest : struct to hold Update Consolidated Shipment request data
type updateRequest struct {
	ID                int64            `json:"-"`
	DriverName        string           `json:"driver_name" valid:"required|alpha_num_space|lte:100"`
	VehicleNumber     string           `json:"vehicle_number" valid:"required|alpha_num|range:5,9"`
	DriverPhoneNumber string           `json:"driver_phone_number" valid:"required|numeric|range:8,15"`
	PurchaseOrders    []*purchaseOrder `json:"purchase_orders" valid:"required"`

	ConsolidatedShipment *model.ConsolidatedShipment `json:"-"`
	Session              *auth.SessionData           `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var err error
	var warehouse, supplierOrganization []string

	c.ConsolidatedShipment, err = repository.ValidConsolidatedShipment(c.ID)
	if err != nil {
		o.Failure("conaolidated_shipment_id.invalid", util.ErrorInvalidData("consolidated shipment id"))
	}

	if c.ConsolidatedShipment.Status != 1 {
		o.Failure("consolidated_shipment.invalid", util.ErrorActive("consolidated shipment"))
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

	c.DriverPhoneNumber = util.ParsePhoneNumberPrefix(c.DriverPhoneNumber)

	return o
}

func (c *updateRequest) Messages() map[string]string {
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
