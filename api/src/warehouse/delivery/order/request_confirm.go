// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"fmt"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type confirmRequest struct {
	ID                 int64                `json:"-"`
	Note               string               `json:"note"`
	DeliveryOrderItems []*itemRequest       `json:"delivery_order_items" valid:"required"`
	CodeDR             string               `json:"-"`
	Session            *auth.SessionData    `json:"-"`
	DeliveryOrder      *model.DeliveryOrder `json:"-"`

	// FOR CREATE DELIVERY RETURN
	DeliveryReturn     *model.DeliveryReturn
	DeliveryReturnItem []*model.DeliveryReturnItem
}

// Validate : function to validate uom request data
func (c *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var deliveryOrderItemID int64

	duplicateProductList := make(map[string]bool)
	duplicate := make(map[string]bool)

	c.DeliveryOrder = &model.DeliveryOrder{ID: c.ID}
	if e = c.DeliveryOrder.Read("ID"); e == nil {
		if c.DeliveryOrder.Status != 1 && c.DeliveryOrder.Status != 6 && c.DeliveryOrder.Status != 7 {
			o.Failure("id.invalid", util.ErrorActive("delivery order"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("delivery order"))
	}

	if e = c.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Failure("sales_order.id.invalid", util.ErrorInvalidData("sales order"))
	}

	if e = c.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
		o.Failure("order_type.id.invalid", util.ErrorInvalidData("order type"))
	}

	for n, row := range c.DeliveryOrderItems {
		var productID int64

		if row.ID != "" {

			if !duplicate[row.ID] {
				deliveryOrderItemID, _ = common.Decrypt(row.ID)
				row.DeliveryOrderItem = &model.DeliveryOrderItem{ID: deliveryOrderItemID}
				if e = row.DeliveryOrderItem.Read("ID"); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
				}
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		}

		if len(row.Note) > 100 {
			o.Failure("note"+strconv.Itoa(n), util.ErrorCharLength("note", 100))
		}

		if !duplicateProductList[row.ProductID] {
			if row.DeliverQty < 0 {
				o.Failure(fmt.Sprintf("items.%d.deliver_qty.required", n), "Qty cant be 0")
			}

			productID, _ = common.Decrypt(row.ProductID)
			row.Product = &model.Product{ID: productID}

			if e = row.Product.Read("ID"); e != nil {
				o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
			}
			duplicateProductList[row.ProductID] = true

			if row.DeliverQty < row.OrderQty {
				var productDeliveryReturn *model.DeliveryReturnItem

				if c.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
					productDeliveryReturn = &model.DeliveryReturnItem{
						Product:           row.Product,
						DeliveryOrderItem: row.DeliveryOrderItem,
						ReturnGoodQty:     0,
						ReturnWasteQty:    row.OrderQty - row.DeliverQty,
					}
				} else {
					productDeliveryReturn = &model.DeliveryReturnItem{
						Product:           row.Product,
						DeliveryOrderItem: row.DeliveryOrderItem,
						ReturnGoodQty:     row.OrderQty - row.DeliverQty,
						ReturnWasteQty:    0,
					}
				}

				c.DeliveryReturnItem = append(c.DeliveryReturnItem, productDeliveryReturn)
			}
		}

	}

	// GET BRANCH CODE FOR DR CODE & NOTIFICATION
	if e = c.DeliveryOrder.SalesOrder.Branch.Read("ID"); e != nil {
		o.Failure("branch.id.invalid", util.ErrorInvalidData("branch"))
	}
	if e = c.DeliveryOrder.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
		o.Failure("merchant.id.invalid", util.ErrorInvalidData("branch"))
	}
	if e = c.DeliveryOrder.SalesOrder.Branch.Merchant.UserMerchant.Read("ID"); e != nil {
		o.Failure("user_merchant.id.invalid", util.ErrorInvalidData("user merchant"))
	}

	if len(c.DeliveryReturnItem) > 0 {
		c.DeliveryReturn = &model.DeliveryReturn{
			RecognitionDate: time.Now(),
			Warehouse:       c.DeliveryOrder.Warehouse,
			DeliveryOrder:   c.DeliveryOrder,
			Status:          1,
		}

	}

	return o
}

// Messages : function to return error validation messages
func (c *confirmRequest) Messages() map[string]string {
	return map[string]string{
		"sales_order_id.required": util.ErrorInputRequired("sales order"),
	}
}
