// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer_sku

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// confirmRequest : struct to hold Confirm Transfer SKU request data
type confirmRequest struct {
	ID                      int64              `json:"-"`
	TransferSku             *model.TransferSku `json:"-"`
	MapAvailableStock       map[int64]float64
	AmountOfChildProduct    map[int64]int
	TotalQtyTransferProduct map[int64]float64
	IsTransferToAnotherSku  map[int64]bool
	TotalQtyWasteProduct    map[int64]float64

	GrTransferSKU bool
	Session       *auth.SessionData `json:"-"`
}

// Validate : function to validate Confirm Transfer SKU request data
func (c *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var filter, exclude map[string]interface{}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	if c.TransferSku, err = repository.GetTransferSku("ID", c.ID); err != nil {
		o.Failure("id_invalid", util.ErrorInvalidData("transfer sku"))
		return o
	}

	if c.TransferSku.Status != 1 {
		o.Failure("status_invalid", util.ErrorActive("transfer sku"))
		return o
	}

	if err = c.TransferSku.Warehouse.Read("ID"); err != nil {
		o.Failure("id_invalid", util.ErrorInvalidData("warehouse"))
		return o
	}

	if c.TransferSku.GoodsReceipt != nil {
		if err = c.TransferSku.GoodsReceipt.Read("ID"); err != nil {
			o.Failure("id_invalid", util.ErrorInvalidData("good receipt"))
			return o
		}
		c.GrTransferSKU = true

		if c.TransferSku.PurchaseOrder != nil {
			if err = c.TransferSku.PurchaseOrder.Read("ID"); err != nil {
				o.Failure("id_invalid", util.ErrorInvalidData("purchase order"))
				return o
			}
		}
	}

	if c.TransferSku.TotalTransferQty > 0 {
		stockType, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "good stock")
		if e != nil {
			o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
			return o
		}
		filter = map[string]interface{}{"status": 1, "warehouse_id": c.TransferSku.Warehouse.ID, "stock_type": stockType.ValueInt}
		exclude = map[string]interface{}{}
		if _, countStockOpname, e := repository.CheckStockOpnameData(filter, exclude); e == nil && countStockOpname > 0 {
			o.Failure("id.invalid", util.ErrorRelated("active ", "stock opname", c.TransferSku.Warehouse.Code+"-"+c.TransferSku.Warehouse.Name))
		}
	}
	if c.TransferSku.TotalWasteQty > 0 {
		stockType, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "waste stock")
		if e != nil {
			o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
			return o
		}
		filter = map[string]interface{}{"status": 1, "warehouse_id": c.TransferSku.Warehouse.ID, "stock_type": stockType.ValueInt}
		exclude = map[string]interface{}{}
		if _, countStockOpname, e := repository.CheckStockOpnameData(filter, exclude); e == nil && countStockOpname > 0 {
			o.Failure("id.invalid", util.ErrorRelated("active ", "stock opname", c.TransferSku.Warehouse.Code+"-"+c.TransferSku.Warehouse.Name))
		}
	}

	c.MapAvailableStock = make(map[int64]float64)
	c.IsTransferToAnotherSku = make(map[int64]bool)
	c.TotalQtyTransferProduct = make(map[int64]float64)
	c.TotalQtyWasteProduct = make(map[int64]float64)

	for _, v := range c.TransferSku.TransferSkuItems {
		if _, ok := c.TotalQtyTransferProduct[v.Product.ID]; ok {
			c.TotalQtyTransferProduct[v.Product.ID] += v.TransferQty * v.TransferProduct.UnitWeight
			c.IsTransferToAnotherSku[v.Product.ID] = true
		} else {
			c.TotalQtyTransferProduct[v.Product.ID] = 0
			c.TotalQtyWasteProduct[v.Product.ID] = v.WasteQty
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
