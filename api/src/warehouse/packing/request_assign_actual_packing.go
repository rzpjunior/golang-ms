// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// updateRequest
type assignActualPackingRequest struct {
	ID                 int64   `json:"-" valid:"required"`
	PackingOrderItemID string  `json:"packing_order_item_id" valid:"required"`
	ProductID          string  `json:"product_id" valid:"required"`
	HelperID           string  `json:"helper_id" valid:"required"`
	TotalPack          float64 `json:"total_pack" valid:"gte:0"`
	TotalWeight        float64 `json:"total_weight" valid:"gte:0"`

	PackingOrderItem *model.PackingOrderItem
	Product          *model.Product
	Helper           *model.Staff

	PackingOrder *model.PackingOrder
	Session      *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *assignActualPackingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.PackingOrder, err = repository.ValidPackingOrder(c.ID); err == nil {
		if c.PackingOrder.Status != 1 {
			o.Failure("packing_order_id.status", util.ErrorDocStatus("packing order", "active"))
		}
	} else {
		o.Failure("packing_order_id.invalid", util.ErrorInvalidData("id"))
	}

	if pID, e := common.Decrypt(c.ProductID); e != nil {
		o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
	} else {
		if c.Product, e = repository.ValidProduct(pID); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}
		if c.Product.Packability != 1 {
			o.Failure("product_id", "Product must be packable")
		}
	}

	if hID, e := common.Decrypt(c.HelperID); e != nil {
		o.Failure("helper_id.invalid", util.ErrorInvalidData("helper"))
	} else {
		if c.Helper, e = repository.ValidStaff(hID); e != nil {
			o.Failure("helper_id.invalid", util.ErrorInvalidData("helper"))
		}
	}

	if packingOrderItemID, e := common.Decrypt(c.PackingOrderItemID); e != nil {
		o.Failure("packing_order_item_id.invalid", util.ErrorInvalidData("packing order item"))
	} else {
		if c.PackingOrderItem, e = repository.ValidPackingOrderItem(packingOrderItemID); e != nil {
			o.Failure("packing_order_item_id.invalid", util.ErrorInvalidData("packing order item"))
		}
	}

	if c.PackingOrderItem.PackingOrder.ID != c.PackingOrder.ID {
		o.Failure("packing_order_id", util.ErrorMustBeSame("packing order id", "packing order id"))
	}


	return o
}

// Messages : function to return error validation messages
func (c *assignActualPackingRequest) Messages() map[string]string {
	messages := map[string]string{
		"packing_order_item_id.required": util.ErrorInputRequired("packing order item id"),
		"product_id.required": util.ErrorInputRequired("product"),
		"helper_id.required": util.ErrorInputRequired("helper"),
	}

	return messages
}
