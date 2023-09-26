// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package disposal

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type confirmRequest struct {
	ID            int64                `json:"-"`
	Session       *auth.SessionData    `json:"-"`
	WasteDisposal *model.WasteDisposal `json:"-"`
	Stocks        []*model.Stock       `json:"-"`
}

// Validate : function to validate uom request data
func (c *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var st model.Stock
	var stockOpname int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	c.WasteDisposal, e = repository.GetWasteDisposal("ID", c.ID)
	if e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("waste disposal"))
	} else {
		if c.WasteDisposal.Status != 1 {
			o.Failure("id.invalid", util.ErrorActive("waste disposal"))
			return o
		}
	}

	if err := c.WasteDisposal.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	stockType, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "waste stock")
	if e != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	orSelect.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND stock_type = ? AND status = 1", c.WasteDisposal.Warehouse.ID, stockType.ValueInt).QueryRow(&stockOpname)

	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", c.WasteDisposal.Warehouse.Name))

	}

	if len(c.WasteDisposal.WasteDisposalItems) > 0 {
		for _, v := range c.WasteDisposal.WasteDisposalItems {
			orSelect.Raw("select * from stock s where s.product_id = ? and s.warehouse_id = ?", v.Product.ID, v.WasteDisposal.Warehouse.ID).QueryRow(&st)
			if st.WasteStock < v.DisposeQty {
				o.Failure("dispose_qty.invalid", util.ErrorEqualLess("dispose qty", "waste stock"))
			}
			stock := &model.Stock{Product: v.Product, Warehouse: c.WasteDisposal.Warehouse}
			stock.Read("Product", "Warehouse")
			c.Stocks = append(c.Stocks, stock)
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
