// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package _return

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold price set request data
type createRequest struct {
	Code                string         `json:"-"`
	RecognitionDate     string         `json:"recognition_date" valid:"required"`
	AreaID              string         `json:"area_id" valid:"required"`
	WarehouseID         string         `json:"warehouse_id" valid:"required"`
	DeliveryOrderID     string         `json:"delivery_order_id" valid:"required"`
	DeliveryReturnItems []*itemRequest `json:"delivery_return_items" valid:"required"`

	Note              string    `json:"note"`
	RecognitionDateAt time.Time `json:"-"`

	DeliveryOrder *model.DeliveryOrder `json:"-"`
	SalesOrder    *model.SalesOrder    `json:"-"`
	Branch        *model.Branch        `json:"-"`
	Warehouse     *model.Warehouse     `json:"-"`
	Area          *model.Area          `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

type itemRequest struct {
	ID                  string  `json:"id"`
	ProductID           string  `json:"product_id"`
	DeliveryOrderItemID string  `json:"delivery_order_item_id"`
	Uom                 string  `json:"uom"`
	DeliverQty          float64 `json:"deliver_qty"`
	ReceiveQty          float64 `json:"receive_qty"`
	ReturnGoodStockQty  float64 `json:"return_good_stock_qty" valid:"required"`
	ReturnWasteStockQty float64 `json:"return_waste_stock_qty" valid:"required"`
	WasteReason         int8    `json:"return_waste_reason"`
	Note                string  `json:"note" valid:"lte:255"`

	DeliveryReturnItem *model.DeliveryReturnItem `json:"-"`
	DeliveryOrderItem  *model.DeliveryOrderItem  `json:"-"`
	Product            *model.Product            `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	layout := "2006-01-02"
	var e error
	var filter, exclude map[string]interface{}

	if c.RecognitionDateAt, e = time.Parse(layout, c.RecognitionDate); e != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("delivery return date"))
	}

	if deliveryOrderID, e := common.Decrypt(c.DeliveryOrderID); e != nil {
		o.Failure("delivery_order_id.invalid", util.ErrorInvalidData("delivery order"))
	} else {
		if c.DeliveryOrder, e = repository.ValidDeliveryOrder(deliveryOrderID); e == nil {
			if c.SalesOrder, e = repository.ValidSalesOrder(c.DeliveryOrder.SalesOrder.ID); e == nil {
				c.Branch, e = repository.ValidBranch(c.SalesOrder.Branch.ID)
			}
		} else {
			o.Failure("delivery_order_id.invalid", util.ErrorInvalidData("delivery order"))
		}
	}

	if whID, e := common.Decrypt(c.WarehouseID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	} else {
		if c.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}
	}

	if areaID, e := common.Decrypt(c.AreaID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	} else {
		if c.Area, e = repository.ValidArea(areaID); e != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		}
	}

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	if e = c.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
	}

	if e = c.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
		o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
	}

	var duplicated = make(map[string]bool)

	var returnGoodStock float64
	var returnWasteStock float64
	for n, v := range c.DeliveryReturnItems {

		if deliveryOrderItemID, e := common.Decrypt(v.DeliveryOrderItemID); e != nil {
			o.Failure("delivery_order_item_id.invalid", util.ErrorInvalidData("delivery order item"))
		} else {
			if v.DeliveryOrderItem, e = repository.ValidDeliveryOrderItem(deliveryOrderItemID); e != nil {
				o.Failure("delivery_order_item_id.invalid", util.ErrorInvalidData("delivery order item"))
			}
		}

		var productID int64
		if v.ProductID != "" {
			if !duplicated[v.ProductID] {
				productID, _ = common.Decrypt(v.ProductID)
				if v.Product, e = repository.ValidProduct(productID); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
				} else {
					filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.Warehouse.ID, "status": 1}
					if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
						o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
					}
				}
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		} else {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}

		if c.DeliveryOrder.SalesOrder.OrderType.Name != "Zero Waste" {
			if v.ReturnWasteStockQty > 0 {
				if v.WasteReason == 0 {
					o.Failure("waste_reason_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("waste reason"))
				} else {
					_, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "waste_reason", "value_int", v.WasteReason)
					if e != nil {
						o.Failure("waste_reason"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("waste_reason"))
					}
				}
			} else {
				v.WasteReason = 0
			}
		} else {
			v.WasteReason = 0
		}

		returnGoodStock += v.ReturnGoodStockQty
		returnWasteStock += v.ReturnWasteStockQty
	}
	if returnGoodStock < 0 {
		o.Failure("id.invalid", util.ErrorEqualGreater("return good stock qty", "0"))
	}

	if returnWasteStock < 0 {
		o.Failure("id.invalid", util.ErrorEqualGreater("return waste stock qty", "0"))
	}

	if returnGoodStock+returnWasteStock == 0 {
		o.Failure("id.invalid", util.ErrorReturnStockCannot0Qty())
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"recognition_date.required":  util.ErrorInputRequired("recognition date"),
		"area_id.required":           util.ErrorInputRequired("area"),
		"warehouse_id.required":      util.ErrorInputRequired("warehouse"),
		"delivery_order_id.required": util.ErrorInputRequired("delivery order"),
	}

	for i := range c.DeliveryReturnItems {
		messages["item."+strconv.Itoa(i)+".product_id.required"] = util.ErrorInputRequired("product")
		messages["item."+strconv.Itoa(i)+".delivery_order_item_id.required"] = util.ErrorInputRequired("purchase order item")
		messages["item."+strconv.Itoa(i)+".return_good_stock_qty.required"] = util.ErrorInputRequired("return good stock qty")
		messages["item."+strconv.Itoa(i)+".return_waste_stock_qty.required"] = util.ErrorInputRequired("return waste stock qty")

	}

	return messages
}
