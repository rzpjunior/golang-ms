// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strconv"
)

// updateRequest : struct to hold price set request data
type updateRequest struct {
	ID int64 `json:"-" valid:"required"`

	PackingOrderItems []*itemRequests `json:"packing_order_items" valid:"required"`
	PackingOrder      *model.PackingOrder
	Session           *auth.SessionData
}

type itemRequests struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`

	TotalOrder  float64  `json:"total_order"`
	TotalWeight float64  `json:"total_weight"`
	TotalPack   float64  `json:"total_pack"`
	Helper      []string `json:"helper" valid:"required"`

	HelperDec []int64 `json:"-"`

	Product *model.Product `json:"-"`
}

// Validate : function to validate uom request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	c.PackingOrder = &model.PackingOrder{ID: c.ID}
	c.PackingOrder.Read("ID")
	if c.PackingOrder.Status != 1 {
		o.Failure("status.inactive", util.ErrorActive("sales order"))
	}

	var duplicated = make(map[string]bool)

	for n, row := range c.PackingOrderItems {
		var productID int64

		if row.ProductID != "" {
			if !duplicated[row.ProductID] {

				productID, _ = common.Decrypt(row.ProductID)
				row.Product = &model.Product{ID: productID}

				if e = row.Product.Read("ID"); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
				}
				duplicated[row.ProductID] = true
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		} else {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}

	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":        util.ErrorInputRequired("warehouse"),
		"area_id.required":             util.ErrorInputRequired("area"),
		"delivery_date.required":       util.ErrorInputRequired("delivery date"),
		"packing_order_items.required": util.ErrorInputRequired("packing order item"),
	}

	return messages
}
