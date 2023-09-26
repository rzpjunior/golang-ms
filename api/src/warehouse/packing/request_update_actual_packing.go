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
	"strconv"
)

// updateRequest
type updateActualPackingRequest struct {
	ID                     int64                `json:"-" valid:"required"`
	PackingOrderItemUpdate []*requestActualItem `json:"packing_update" valid:"required"`

	//Warehouse *model.Warehouse
	PackingOrder *model.PackingOrder
	Session      *auth.SessionData `json:"-"`
}

type requestActualItem struct {
	PackingOrderItemID string  `json:"packing_order_item_id" valid:"required"`
	ProductID          string  `json:"product_id" valid:"required"`
	HelperID           string  `json:"helper_id"`
	TotalPack          float64 `json:"total_pack" valid:"gte:0"`
	TotalWeight        float64 `json:"total_weight" valid:"gte:0"`

	PackingOrderItem *model.PackingOrderItem
	Product          *model.Product
	Helper           *model.Staff
}

// Validate : function to validate uom request data
func (c *updateActualPackingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.PackingOrder, err = repository.ValidPackingOrder(c.ID); err == nil {
		if c.PackingOrder.Status != 1 {
			o.Failure("upload_packing_order.error", util.ErrorDocStatus("packing order", "active"))
		}
	} else {
		o.Failure("upload_packing_order.error", util.ErrorInvalidData("id"))
	}

	for _, v := range c.PackingOrderItemUpdate {
		if pID, e := common.Decrypt(v.ProductID); e != nil {
			o.Failure("upload_packing_order.error", util.ErrorInvalidData("product"))
		} else {
			if v.Product, e = repository.ValidProduct(pID); e != nil {
				o.Failure("upload_packing_order.error", util.ErrorInvalidData("product"))
			}
			if v.Product.Packability != 1 {
				o.Failure("upload_packing_order.error", "Product must be packable")
			}
		}

		if v.HelperID != "" {
			if hID, e := common.Decrypt(v.HelperID); e != nil {
				o.Failure("upload_packing_order.error", util.ErrorInvalidData("helper"))
			} else {
				if v.Helper, e = repository.ValidStaff(hID); e != nil {
					o.Failure("upload_packing_order.error", util.ErrorInvalidData("helper"))
				}
			}
		}

		if packingOrderItemID, e := common.Decrypt(v.PackingOrderItemID); e != nil {
			o.Failure("upload_packing_order.error", util.ErrorInvalidData("packing order item"))
		} else {
			if v.PackingOrderItem, e = repository.ValidPackingOrderItem(packingOrderItemID); e != nil {
				o.Failure("upload_packing_order.error", util.ErrorInvalidData("packing order item"))
			}
		}

		if v.PackingOrderItem.PackingOrder.ID != c.PackingOrder.ID {
			o.Failure("upload_packing_order.error", util.ErrorMustBeSame("code", "packing order code"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateActualPackingRequest) Messages() map[string]string {
	messages := map[string]string{
		"packing_update.required": util.ErrorInputRequired("product"),
	}

	for i, _ := range c.PackingOrderItemUpdate {
		messages["upload_packing_order."+strconv.Itoa(i)+".required"] = util.ErrorInputRequired("packing order item")
		messages["upload_packing_order."+strconv.Itoa(i)+".required"] = util.ErrorInputRequired("product")
	}

	return messages
}
