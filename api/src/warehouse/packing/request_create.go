// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"time"

	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold price set request data
type createRequest struct {
	Code             string    `json:"-"`
	WarehouseID      string    `json:"warehouse_id" valid:"required"`
	DeliveryDate     string    `json:"delivery_date" valid:"required"`
	Note             string    `json:"note"`
	DeliveryDateTime time.Time `json:"-"`

	Warehouse         *model.Warehouse `json:"-"`
	PackingOrderItems []*itemRequest   `json:"packing_order_items" valid:"required"`
	CollectedSOID     []int64          `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type itemRequest struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`

	TotalOrder   float64  `json:"total_order"`
	TotalWeight  float64  `json:"total_weight"`
	TotalPack    float64  `json:"total_pack"`
	Helper       []string `json:"helper"`
	SalesOrderID string   `json:"sales_order_id"`

	HelperDec []int64 `json:"-"`

	Product    *model.Product    `json:"-"`
	SalesOrder *model.SalesOrder `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var countDoc int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	warehouseID, e := common.Decrypt(c.WarehouseID)
	if e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}
	c.Warehouse = &model.Warehouse{ID: warehouseID}
	if e := c.Warehouse.Read("ID"); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	orSelect.Raw("SELECT COUNT(id) FROM packing_order where delivery_date = ? AND warehouse_id = ? AND status IN (1,2)", c.DeliveryDate, warehouseID).QueryRow(&countDoc)
	if countDoc > 0 {
		o.Failure("delivery_date", util.ErrorInputCannotBeSame("delivery date", "existing delivery date"))
	}

	if c.DeliveryDateTime, e = time.Parse("2006-01-02", c.DeliveryDate); e != nil {
		o.Failure("delivery_date.invalid", util.ErrorInputRequired("delivery date"))
	}

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	var duplicated = make(map[string]bool)

	for n, row := range c.PackingOrderItems {
		var productID int64

		if row.ProductID != "" {
			if !duplicated[row.ProductID] {

				productID, _ = common.Decrypt(row.ProductID)
				row.Product = &model.Product{ID: productID}

				if e = row.Product.Read("ID"); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
				}
				duplicated[row.ProductID] = true
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		} else {
			o.Failure("product_id"+strconv.Itoa(n)+".required", util.ErrorInputRequired("product"))
		}

		if len(row.Helper) > 0 {
			for _, v := range row.Helper {

				helperId, _ := common.Decrypt(v)

				if helperPerson, err := repository.ValidStaff(helperId); err != nil {
					o.Failure("helper"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("helper"))
				} else {
					if helperPerson.Status != int8(1) {
						o.Failure("helper"+strconv.Itoa(n)+".active", util.ErrorActive("helper"))
					}

				}
			}
		}

	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":        util.ErrorInputRequired("warehouse"),
		"delivery_date.required":       util.ErrorInputRequired("delivery date"),
		"packing_order_items.required": util.ErrorInputRequired("packing order item"),
	}

	return messages
}
