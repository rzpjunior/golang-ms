// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createRequest : struct to hold delivery order request data
type createRequest struct {
	JobsID      primitive.ObjectID `json:"jobs_id"`
	WarehouseID string             `json:"warehouse_id" valid:"required"`

	Warehouse          *model.Warehouse  `json:"-"`
	DeliveryOrderItems []*itemRequest    `json:"delivery_order_items" valid:"required"`
	Session            *auth.SessionData `json:"-"`
}

type itemRequest struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Uom         string  `json:"uom"`
	OrderQty    float64 `json:"order_qty"`
	DeliverQty  float64 `json:"deliver_qty"`
	Note        string  `json:"note" valid:"lte:255"`
	Weight      float64 `json:"weight"`
	UnitPrice   float64 `json:"-"`

	DeliveryOrderItem *model.DeliveryOrderItem `json:"-"`
	Product           *model.Product           `json:"-"`
}

// createRequest : function to validate create delivery order based from request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}
	var stockOpname int64

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	warehouseID, e := common.Decrypt(c.WarehouseID)
	if e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("warehouse"))
	}

	c.Warehouse = &model.Warehouse{ID: warehouseID}
	if e = c.Warehouse.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("warehouse"))
	}

	orSelect.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND status = 1", warehouseID).QueryRow(&stockOpname)
	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", c.Warehouse.Name))

	}

	var duplicated = make(map[string]bool)

	for n, row := range c.DeliveryOrderItems {
		var productID int64

		if row.ProductID == "" {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}

		if duplicated[row.ProductID] {
			o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
		}

		if row.DeliverQty < 0 {
			o.Failure("qty"+strconv.Itoa(n)+".greater", util.ErrorGreater("product quantity", "0"))
		}

		productID, _ = common.Decrypt(row.ProductID)
		row.Product = &model.Product{ID: productID}

		if e = row.Product.Read("ID"); e != nil {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
		}

		filter = map[string]interface{}{"product_id": productID, "warehouse_id": warehouseID, "status": 1}

		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
		}

		duplicated[row.ProductID] = true

	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required": util.ErrorInputRequired("warehouse"),
	}

	return messages
}
