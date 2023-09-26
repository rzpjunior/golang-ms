// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type confirmRequest struct {
	ID          int64  `json:"-"`
	InboundType string `json:"inbound_type" valid:"required"`

	GoodsReceipt    *model.GoodsReceipt    `json:"-"`
	Stocks          []*model.Stock         `json:"-"`
	GRQty           []*goodsReceiptQty     `json:"-"`
	PurchaseInvoice *model.PurchaseInvoice `json:"-"`
	StockType       *model.Glossary        `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type goodsReceiptQty struct {
	ReceiveQty   float64
	Note         string
	UnitPrice    float64
	UnitPriceTax float64
}

func (r *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var filter, exclude map[string]interface{}
	var err error

	r.GoodsReceipt, err = repository.GetGoodsReceipt("id", r.ID)
	if err != nil {
		o.Failure("good_receipt.invalid", util.ErrorInvalidData("goods receipt"))
		return o
	}

	if r.GoodsReceipt.Status != 1 {
		o.Failure("status.inactive", util.ErrorDocStatus("goods receipt", "active"))
	}

	r.StockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", r.GoodsReceipt.StockType)
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	filter = map[string]interface{}{"status": 1, "warehouse_id": r.GoodsReceipt.Warehouse.ID, "stock_type": r.StockType.ValueInt}
	exclude = map[string]interface{}{}
	if _, countStockOpname, e := repository.CheckStockOpnameData(filter, exclude); e == nil && countStockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active ", "stock opname", r.GoodsReceipt.Warehouse.Code+"-"+r.GoodsReceipt.Warehouse.Name))
	}

	for _, v := range r.GoodsReceipt.GoodsReceiptItems {
		if v.ReceiveQty > 0 {
			stock := &model.Stock{Product: v.Product, Warehouse: r.GoodsReceipt.Warehouse}
			if e := stock.Read("Product", "Warehouse"); e != nil {
				o.Failure("id_invalid", util.ErrorInvalidData("stock"))
			}
			r.Stocks = append(r.Stocks, stock)
			switch r.InboundType {
			case "goods_transfer":
				grQty := &goodsReceiptQty{ReceiveQty: v.ReceiveQty, Note: v.Note}
				r.GRQty = append(r.GRQty, grQty)
			case "purchase_order":
				if e := v.PurchaseOrderItem.Read("ID"); e != nil {
					o.Failure("id_invalid", util.ErrorInvalidData("purchase order item"))
				}
				grQty := &goodsReceiptQty{ReceiveQty: v.ReceiveQty, Note: v.Note, UnitPrice: v.PurchaseOrderItem.UnitPrice, UnitPriceTax: v.PurchaseOrderItem.UnitPriceTax}
				r.GRQty = append(r.GRQty, grQty)
			}

		}
	}

	switch r.InboundType {
	case "goods_transfer":
		if e := r.GoodsReceipt.GoodsTransfer.Read("ID"); e != nil {
			o.Failure("id_invalid", util.ErrorInvalidData("goods transfer"))
		}
	case "purchase_order":
		if e := r.GoodsReceipt.PurchaseOrder.Read("ID"); e != nil {
			o.Failure("id_invalid", util.ErrorInvalidData("purchase order"))
		}
		r.PurchaseInvoice = &model.PurchaseInvoice{
			PurchaseOrder: r.GoodsReceipt.PurchaseOrder,
		}

		r.PurchaseInvoice.Read("PurchaseOrder")
	default:
		o.Failure("inbound_type", util.ErrorInvalidData("inbound type"))
	}

	return o
}

func (r *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
