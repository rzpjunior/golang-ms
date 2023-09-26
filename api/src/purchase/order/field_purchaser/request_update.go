// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field_purchaser

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequest : struct to hold Update Field Purchase Order request data
type updateRequest struct {
	ID                      int64   `json:"-"`
	PaymentMethodID         string  `json:"payment_method_id" valid:"required"`
	FieldPurchaseOrderItems []*item `json:"field_purchase_order_items" valid:"required"`

	TotalPrice         float64                   `json:"-"`
	TotalItem          int8                      `json:"-"`
	PurchaseOrder      *model.PurchaseOrder      `json:"-"`
	FieldPurchaseOrder *model.FieldPurchaseOrder `json:"-"`
	PurchaseDeliver    *model.PurchaseDeliver    `json:"-"`
	PaymentMethod      *model.PaymentMethod      `json:"-"`
	Session            *auth.SessionData         `json:"-"`
}

type item struct {
	FieldPurchaseOrderItemID string  `json:"field_purchase_order_item_id" valid:"required"`
	PurchaseQty              float64 `json:"purchase_qty" valid:"required|lte:99999999"`
	UnitPrice                float64 `json:"unit_price" valid:"required|lte:9999999999"`

	Product                *model.Product                `json:"-"`
	PurchaseOrderItem      *model.PurchaseOrderItem      `json:"-"`
	FieldPurchaseOrderItem *model.FieldPurchaseOrderItem `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	c.FieldPurchaseOrder = &model.FieldPurchaseOrder{ID: c.ID}
	if err := c.FieldPurchaseOrder.Read("ID"); err != nil {
		o.Failure("field_purchase_order_id.invalid", util.ErrorInvalidData("field purchase order"))
	}

	paymentMethodID, err := common.Decrypt(c.PaymentMethodID)
	if err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
	}

	c.PaymentMethod, err = repository.ValidPaymentMethod(paymentMethodID)
	if err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
	}

	c.PurchaseOrder, err = repository.ValidPurchaseOrder(c.FieldPurchaseOrder.PurchaseOrder)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	c.PurchaseDeliver = &model.PurchaseDeliver{FieldPurchaseOrder: c.FieldPurchaseOrder}
	if err := c.PurchaseDeliver.Read("FieldPurchaseOrder"); err != nil {
		o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver"))
	}

	if c.PurchaseDeliver.DeltaPrint > 1 {
		o.Failure("purchase_deliver.invalid", util.ErrorCannotUpdateAfter("Field Purchase Order", "print Surat Jalan"))
	}

	for i, v := range c.FieldPurchaseOrderItems {

		fieldPurchaseOrderItemID, err := common.Decrypt(v.FieldPurchaseOrderItemID)
		if err != nil {
			o.Failure("field_purchase_order_item_id.invalid", util.ErrorInvalidData("field purchase order item"))
		}

		if v.FieldPurchaseOrderItem, err = repository.ValidFieldPurchaseOrderItem(fieldPurchaseOrderItemID); err != nil {
			o.Failure("field_purchase_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("field purchase order item"))
		}

		if v.PurchaseOrderItem, err = repository.ValidPurchaseOrderItem(v.FieldPurchaseOrderItem.PurchaseOrderItem.ID); err != nil {
			o.Failure("purchase_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase order item"))
		}

		c.PurchaseOrder.TotalPrice -= v.PurchaseOrderItem.Subtotal
		c.PurchaseOrder.TaxAmount -= v.PurchaseOrderItem.TaxAmount
		c.TotalPrice += (v.PurchaseQty * v.UnitPrice)
		c.TotalItem += 1
		v.PurchaseOrderItem.PurchaseQty += v.PurchaseQty

		//count average unit price
		filter := map[string]interface{}{"purchase_order_item_id": v.FieldPurchaseOrderItem.PurchaseOrderItem.ID}
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

func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"payment_method_id.required":          util.ErrorSelectRequired("payment method"),
		"field_purchase_order_items.required": util.ErrorInputRequired("product item"),
	}

	for i, _ := range c.FieldPurchaseOrderItems {
		messages["item."+strconv.Itoa(i)+".field_purchase_order_item_id.required"] = util.ErrorInputRequired("product")
		messages["item."+strconv.Itoa(i)+".purchase_qty.required"] = util.ErrorInputRequired("purchase qty")
		messages["item."+strconv.Itoa(i)+".unit_price.required"] = util.ErrorInputRequired("unit price")
		messages["item."+strconv.Itoa(i)+".purchase_qty.lte"] = util.ErrorEqualLess("purchase qty", "99999999")
		messages["item."+strconv.Itoa(i)+".unit_price.lte"] = util.ErrorEqualLess("unit price", "9999999999")
	}

	return messages
}
