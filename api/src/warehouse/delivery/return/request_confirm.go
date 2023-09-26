// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package _return

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type confirmRequest struct {
	ID             int64                 `json:"-"`
	DeliveryReturn *model.DeliveryReturn `json:"-"`

	StocksAva   map[int]*model.Stock `json:"-"`
	StocksWaste map[int]*model.Stock `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	//var e error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var countOpname int
	var err error
	c.DeliveryReturn, _ = repository.GetDeliveryReturn("ID", c.ID)
	if c.DeliveryReturn.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("delivery return"))
		return o
	}

	if err = c.DeliveryReturn.DeliveryOrder.Read("ID"); err != nil {
		o.Failure("delivery_order_id.invalid", util.ErrorInvalidData("delivery order"))
	}
	if err = c.DeliveryReturn.DeliveryOrder.SalesOrder.Read("ID"); err != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
	}
	if err = c.DeliveryReturn.DeliveryOrder.SalesOrder.OrderType.Read("ID"); err != nil {
		o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
	}

	wasteStockType, err := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "waste stock")
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}
	goodStockType, err := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "good stock")
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	c.StocksAva = make(map[int]*model.Stock)
	c.StocksWaste = make(map[int]*model.Stock)

	for i, v := range c.DeliveryReturn.DeliveryReturnItems {
		if v.ReturnGoodQty > 0 {
			stock := &model.Stock{Product: v.Product, Warehouse: c.DeliveryReturn.Warehouse}
			stock.Read("Product", "Warehouse")
			c.StocksAva[i] = stock

			orSelect.Raw("select count(*) from stock_opname so where so.warehouse_id = ? and so.stock_type = ? and status = 1", c.DeliveryReturn.Warehouse.ID, goodStockType.ValueInt).QueryRow(&countOpname)
			if countOpname > 0 {
				o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", "warehouse"))
				return o
			}
		}

		if v.ReturnWasteQty > 0 {
			stock := &model.Stock{Product: v.Product, Warehouse: c.DeliveryReturn.Warehouse}
			stock.Read("Product", "Warehouse")
			c.StocksWaste[i] = stock

			orSelect.Raw("select count(*) from stock_opname so where so.warehouse_id = ? and so.stock_type = ? and status = 1", c.DeliveryReturn.Warehouse.ID, wasteStockType.ValueInt).QueryRow(&countOpname)
			if countOpname > 0 {
				o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", "warehouse"))
				return o
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
