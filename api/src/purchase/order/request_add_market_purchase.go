// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type marketPurchaseRequest struct {
	ID                   int64                      `json:"-"`
	PurchaseOrderItems   []*requestItem             `json:"purchase_order_items" valid:"required"`
	PurchaseOrder        *model.PurchaseOrder       `json:"-"`
	PurchaseOrderItemArr []*model.PurchaseOrderItem `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (c *marketPurchaseRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	var (
		err               error
		marketPurchase    []byte
		totalCharge       float64
		totalPrice        float64
		marketPurchaseStr string
	)

	c.PurchaseOrder = &model.PurchaseOrder{ID: c.ID}
	if err = c.PurchaseOrder.Read("ID"); err == nil {
		if c.PurchaseOrder.Status != 5 {
			o.Failure("id.invalid", util.ErrorDraft("purchase order"))
			return o
		}
	}

	for _, v := range c.PurchaseOrderItems {
		var (
			poiID int64
		)

		if marketPurchase, err = json.Marshal(v.MarketPurchase); err == nil {
			marketPurchaseStr = string(marketPurchase)
		} else {
			o.Failure("id.invalid", util.ErrorInvalidData("market purchase data"))
		}

		poiID, _ = common.Decrypt(v.ID)
		v.PurchaseOrderItem = &model.PurchaseOrderItem{ID: poiID}
		if err = v.PurchaseOrderItem.Read("ID"); err == nil {
			v.PurchaseOrderItem.PurchaseQty = v.PurchaseQty
			v.PurchaseOrderItem.UnitPrice = v.UnitPrice
			v.PurchaseOrderItem.Subtotal = v.PurchaseQty * v.UnitPrice
			v.PurchaseOrderItem.MarketPurchaseStr = marketPurchaseStr

			c.PurchaseOrderItemArr = append(c.PurchaseOrderItemArr, v.PurchaseOrderItem)

			totalPrice = totalPrice + v.PurchaseOrderItem.Subtotal
		}
	}

	totalCharge = totalPrice + c.PurchaseOrder.DeliveryFee + (c.PurchaseOrder.TaxPct * totalPrice / 100)

	c.PurchaseOrder.TotalPrice = totalPrice
	c.PurchaseOrder.TotalCharge = common.Rounder(totalCharge, 0.5, 2)

	return o
}

func (c *marketPurchaseRequest) Messages() map[string]string {
	return map[string]string{
		"purchase_order_items.required": util.ErrorInputRequired("market purchase"),
	}
}
