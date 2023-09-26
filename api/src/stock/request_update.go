// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type updateRequest struct {
	JobsID primitive.ObjectID `json:"jobs_id"`

	Warehouse          *model.Warehouse     `json:"-"`
	DeliveryOrder      *model.DeliveryOrder `json:"-"`
	DeliveryOrderItems []*itemRequest       `json:"delivery_order_items" valid:"required"`
	Session            *auth.SessionData    `json:"-"`
}

// updateRequest : function to validate update delivery order based from request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var e error
	var filter, exclude map[string]interface{}
	var duplicated = make(map[string]bool)
	var duplicatedoi = make(map[string]bool)
	var deliveryOrderItemID, productID, stockOpname int64

	c.DeliveryOrder = &model.DeliveryOrder{ID: c.DeliveryOrder.ID}
	if e = c.DeliveryOrder.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("delivery order"))
	}

	if c.DeliveryOrder.Status != 1 && c.DeliveryOrder.Status != 5 && c.DeliveryOrder.Status != 6 && c.DeliveryOrder.Status != 7 {
		o.Failure("id.invalid", util.ErrorActive("delivery order"))
		return o
	}

	// DIRECT INVOICE SHOULD HAVE AT LEAST 1 INVOICE - [INVOICE DOC CREATED WHEN SO CREATED]
	o1.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND status = 1", c.DeliveryOrder.Warehouse.ID).QueryRow(&stockOpname)
	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", c.Warehouse.Name))

	}

	for n, row := range c.DeliveryOrderItems {
		if row.ID == "" {
			o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
		}
		if row.DeliverQty < 0 {
			o.Failure("qty"+strconv.Itoa(n)+".greater", util.ErrorGreater("product quantity", "0"))
		}
		if !duplicatedoi[row.ID] {
			deliveryOrderItemID, _ = common.Decrypt(row.ID)
			row.DeliveryOrderItem = &model.DeliveryOrderItem{ID: deliveryOrderItemID}

			if e = row.DeliveryOrderItem.Read("ID"); e != nil {
				o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
			}

		}

		if row.ProductID == "" {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}
		if duplicated[row.ProductID] {
			o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
		}

		productID, _ = common.Decrypt(row.ProductID)
		row.Product = &model.Product{ID: productID}

		if e = row.Product.Read("ID"); e != nil {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
		}
		filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.DeliveryOrder.Warehouse.ID, "status": 1}
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
		}
		duplicated[row.ProductID] = true

	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{}
}
