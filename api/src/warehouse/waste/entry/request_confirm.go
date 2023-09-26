// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entry

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type confirmRequest struct {
	ID int64 `json:"-"`

	WasteEntryItems []*items          `json:"-"`
	Session         *auth.SessionData `json:"-"`
	WasteEntry      *model.WasteEntry `json:"-"`
	Warehouse       *model.Warehouse  `json:"-"`
}

// Validate : function to validate uom request data
func (c *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}
	var stockOpname int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	duplicateProductList := make(map[string]bool)

	c.WasteEntry = &model.WasteEntry{ID: c.ID}
	if e = c.WasteEntry.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("waste entry"))
		return o
	}

	if c.WasteEntry.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("waste entry"))
	}

	warehouse, e := repository.ValidWarehouse(c.WasteEntry.Warehouse.ID)
	if e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		return o
	}

	if warehouse.Status != int8(1) {
		o.Failure("warehouse_id.invalid", util.ErrorActive("warehouse"))
	}

	orSelect.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND status = 1", warehouse.ID).QueryRow(&stockOpname)

	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", warehouse.Name))
	}

	orSelect.LoadRelated(c.WasteEntry, "WasteEntryItems", 2)
	for n, row := range c.WasteEntryItems {
		var productID int64

		if !duplicateProductList[row.ProductID] {

			productID, _ = common.Decrypt(row.ProductID)
			row.Product = &model.Product{ID: productID}

			if e = row.Product.Read("ID"); e != nil {
				o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
			}

			filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.WasteEntry.Warehouse.ID, "status": 1}
			if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
				o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
			}

			duplicateProductList[row.ProductID] = true
		}

	}

	return o
}

// Messages : function to return error validation messages
func (c *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
