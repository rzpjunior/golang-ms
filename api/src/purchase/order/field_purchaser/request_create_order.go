// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field_purchaser

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createOrderRequest : struct to hold input create field purchase order request data
type createOrderRequest struct {
	Code                string  `json:"-"`
	CodePurchaseDeliver string  `json:"-"`
	PurchaseOrderID     string  `json:"purchase_order_id" valid:"required"`
	StallID             string  `json:"stall_id" valid:"required"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Signature           string  `json:"signature"`
	Name                string  `json:"name"`
	PaymentMethodID     string  `json:"payment_method_id" valid:"required"`
	Items               []*Item `json:"items" valid:"required"`

	TotalPrice    float64              `json:"-"`
	TotalItem     int8                 `json:"-"`
	PurchaseOrder *model.PurchaseOrder `json:"-"`
	Stall         *model.Stall         `json:"-"`
	PaymentMethod *model.PaymentMethod `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

type Item struct {
	PurchaseOrderItemID string  `json:"purchase_order_item_id" valid:"required"`
	PurchaseQty         float64 `json:"purchase_qty" valid:"required|lte:99999999"`
	UnitPrice           float64 `json:"unit_price" valid:"required|lte:9999999999"`

	Product           *model.Product           `json:"-"`
	PurchaseOrderItem *model.PurchaseOrderItem `json:"-"`
}

// Validate : function to validate create field purchase order request data
func (c *createOrderRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// validate purchase order
	purchaseOrderID, err := common.Decrypt(c.PurchaseOrderID)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order id"))
	}

	c.PurchaseOrder = &model.PurchaseOrder{ID: purchaseOrderID}
	if err = c.PurchaseOrder.Read("ID"); err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	if c.Code, err = util.CheckTable("field_purchase_order"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if c.CodePurchaseDeliver, err = util.CheckTable("purchase_deliver"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	// check if purchase order is assigned to logged in user
	if c.PurchaseOrder.AssignedTo.ID != c.Session.Staff.ID {
		o.Failure("user.forbidden", util.ErrorNotValidFor("purchase order", c.Session.Staff.Name))
	}

	// get warehouse code
	if err = c.PurchaseOrder.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse id"))
	}

	// check if status po is not draft
	if c.PurchaseOrder.Status != 5 {
		o.Failure("purchase_order_id.invalid", util.ErrorDraft("purchase order"))
	}

	// validate stall
	stallID, err := common.Decrypt(c.StallID)
	if err != nil {
		o.Failure("stall_id.invalid", util.ErrorInvalidData("stall id"))
	}

	c.Stall = &model.Stall{ID: stallID}
	if err = c.Stall.Read("ID"); err != nil {
		o.Failure("stall_id.invalid", util.ErrorInvalidData("stall id"))
	}

	// validate payment method
	paymentMethodID, err := common.Decrypt(c.PaymentMethodID)
	if err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method id"))
	}

	c.PaymentMethod = &model.PaymentMethod{ID: paymentMethodID}
	if err = c.PaymentMethod.Read("ID"); err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method id"))
	}

	if c.Name != "" {
		if len(c.Name) > 100 {
			o.Failure("latitude", util.ErrorEqualLess("name", "100 characters"))
		}
	}

	// check creator role
	if !(c.Session.Staff.Role.Name == "Sourcing Admin" || c.Session.Staff.Role.Name == "Field Purchaser") {
		o.Failure("field_purchaser_role.invalid", util.ErrorRole("field purchaser", "sourcing admin or field purchaser"))
	}

	for _, v := range c.Items {
		purchaseOrderItemID, err := common.Decrypt(v.PurchaseOrderItemID)
		if err != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item id"))
		}

		v.PurchaseOrderItem = &model.PurchaseOrderItem{ID: purchaseOrderItemID}
		if err = v.PurchaseOrderItem.Read("ID"); err != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item id"))
		}

		// check if purchase order item is child of purchase order
		if v.PurchaseOrderItem.PurchaseOrder.ID != c.PurchaseOrder.ID {
			o.Failure("purchase_order_item.invalid", util.ErrorNotValidFor("purchase order item", "purchase order"))
		}

		// get product object
		v.Product = &model.Product{ID: v.PurchaseOrderItem.Product.ID}
		if err = v.Product.Read("ID"); err != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product id"))
		}

		c.PurchaseOrder.TotalPrice -= v.PurchaseOrderItem.Subtotal
		c.PurchaseOrder.TaxAmount -= v.PurchaseOrderItem.TaxAmount
		c.TotalPrice += (v.PurchaseQty * v.UnitPrice)
		c.TotalItem += 1
		v.PurchaseOrderItem.PurchaseQty += v.PurchaseQty

		//count average unit price
		filter := map[string]interface{}{"purchase_order_item_id": purchaseOrderItemID}
		exclude := map[string]interface{}{}
		fieldPuchaseOrderItems, _, err := repository.GetFilterFieldPurchaseOrderItems(filter, exclude)
		if err != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item"))
		}

		var subtotal float64
		for _, item := range fieldPuchaseOrderItems {
			subtotal += (item.PurchaseQty * item.UnitPrice)
		}

		// new unit price
		unitPrice := common.Rounder((subtotal+(v.PurchaseQty*v.UnitPrice))/v.PurchaseOrderItem.PurchaseQty, 0.5, 2)

		switch v.PurchaseOrderItem.TaxableItem {
		// condition when product taxable
		case 1:
			switch v.PurchaseOrderItem.IncludeTax {
			// condition when unit price is include tax
			case 1:
				v.PurchaseOrderItem.UnitPrice = unitPrice - (unitPrice * (v.PurchaseOrderItem.TaxPercentage / (100 + v.PurchaseOrderItem.TaxPercentage)))
				v.PurchaseOrderItem.UnitPriceTax = unitPrice
			// condition when unit price is not include tax
			case 2:
				v.PurchaseOrderItem.UnitPriceTax = unitPrice + (unitPrice * (v.PurchaseOrderItem.TaxPercentage / 100))
				v.PurchaseOrderItem.UnitPrice = unitPrice
			}
		// condition when product is non taxable
		case 2:
			v.PurchaseOrderItem.UnitPrice = unitPrice
			v.PurchaseOrderItem.UnitPriceTax = 0
		}

		v.PurchaseOrderItem.Subtotal = v.PurchaseOrderItem.UnitPrice * v.PurchaseOrderItem.PurchaseQty
		v.PurchaseOrderItem.TaxAmount = common.Rounder(v.PurchaseOrderItem.Subtotal*(v.PurchaseOrderItem.TaxPercentage/100), 0.5, 2)
		c.PurchaseOrder.TaxAmount += v.PurchaseOrderItem.TaxAmount
		c.PurchaseOrder.TotalPrice += v.PurchaseOrderItem.Subtotal
	}

	c.PurchaseOrder.TotalCharge = c.PurchaseOrder.TotalPrice + c.PurchaseOrder.TaxAmount + c.PurchaseOrder.DeliveryFee

	return o
}

func (c *createOrderRequest) Messages() map[string]string {
	return map[string]string{
		"stall_id.required":            util.ErrorSelectRequired("stall"),
		"purchase_order_id.required":   util.ErrorInputRequired("purchase order"),
		"items.required":               util.ErrorInputRequired("items"),
		"purchase_order_item.required": util.ErrorInputRequired("purchase order item"),
		"purchase_qty.required":        util.ErrorInputRequired("purchase qty"),
		"unit_price.required":          util.ErrorInputRequired("unit price"),
		"purchase_qty.lte":             util.ErrorEqualLess("purchase qty", "99999999"),
		"unit_price.lte":               util.ErrorEqualLess("unit price", "9999999999"),
		"payment_method_id.required":   util.ErrorSelectRequired("payment method"),
	}
}
