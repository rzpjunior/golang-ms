// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"math"
	"strconv"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateProductRequest struct {
	ID                  int64   `json:"-"`
	DeliveryFee         float64 `json:"delivery_fee"`
	TaxPct              float64 `json:"tax_pct"`
	TotalCharge         float64 `json:"-"`
	TotalPrice          float64 `json:"-"`
	TaxAmount           float64 `json:"-"`
	TotalPriceDebitNote float64 `json:"-"`
	TotalWeight         float64 `json:"-"`

	PurchaseOrderReqItems []*requestItem                 `json:"purchase_order_items" valid:"required"`
	PurchaseOrder         *model.PurchaseOrder           `json:"-"`
	PurchaseOrderItems    []*model.PurchaseOrderItem     `json:"-"`
	DebitNote             *model.DebitNote               `json:"-"`
	MapDebitNoteItem      map[int64]*model.DebitNoteItem `json:"-"`
	Images                []*image                       `json:"images"`
	PurchasePlan          *model.PurchasePlan            `json:"-"`
	Uom                   *model.Uom                     `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type image struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

func (c *updateProductRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q1 := orm.NewOrm()
	q1.Using("read_only")

	var e error
	var productID int64
	var filter, exclude map[string]interface{}
	var gr []*model.GoodsReceipt

	c.PurchaseOrder = &model.PurchaseOrder{ID: c.ID}
	if e = c.PurchaseOrder.Read("ID"); e != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	filter = map[string]interface{}{"purchase_order_id": c.PurchaseOrder.ID, "status__in": []int{1, 2}}
	gr, _, _ = repository.CheckGoodsReceiptData(filter, exclude)

	// region debit note
	q1.Raw("SELECT dn.id, dn.code FROM debit_note dn "+
		"join supplier_return sr on dn.supplier_return_id = sr.id "+
		"join goods_receipt gr on sr.goods_receipt_id = gr.id "+
		"where dn.status = 1 and gr.purchase_order_id = ?", c.ID).QueryRow(&c.DebitNote)

	if c.DebitNote != nil {
		if c.DebitNote.UsedInPurchaseInvoice == 1 {
			o.Failure("debit_note_id.invalid", util.ErrorIsBeingUsed("debit note"))
			return o
		}
		q1.LoadRelated(c.DebitNote, "DebitNoteItems", 1)

		c.MapDebitNoteItem = make(map[int64]*model.DebitNoteItem)
		for _, dni := range c.DebitNote.DebitNoteItems {
			c.MapDebitNoteItem[dni.Product.ID] = dni
		}
	}
	// endregion

	if c.PurchaseOrder.PurchasePlan != nil {
		c.PurchasePlan, e = repository.ValidPurchasePlan(c.PurchaseOrder.PurchasePlan.ID)
		if e != nil {
			o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
		}

		if c.PurchasePlan.Status == 3 {
			o.Failure("purchase_plan_id.invalid", util.ErrorDocStatus("purchase plan", "active or finish"))
		}

		if len(c.Images) == 0 {
			o.Failure("purchase_order_image.required", util.ErrorInputRequired("image"))
		}
	}

	filter = map[string]interface{}{"purchase_order_id": c.ID, "status__in": []int8{1, 2, 6}}
	_, total, e := repository.GetDataPurchaseInvoice(filter, exclude)
	if e != nil {
		o.Failure("purchase_invoice.invalid", util.ErrorInvalidData("purchase invoice"))
	}

	if total > 0 {
		o.Failure("purchase_order.invalid", util.ErrorCannotUpdateAfter("purchase order", "purchase invoice"))
	}

	for i, v := range c.PurchaseOrderReqItems {

		if v.ProductID == "" {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("product"))
		}

		if productID, e = common.Decrypt(v.ProductID); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		v.Product, e = repository.ValidProduct(productID)
		if e != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}

		v.Uom, e = repository.ValidUom(v.Product.Uom.ID)
		if e != nil {
			o.Failure("uom_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("uom"))
		}

		if v.Uom.DecimalEnabled == 2 {
			if math.Mod(v.OrderQty, 1) != 0 {
				o.Failure("order_qty"+strconv.Itoa(i)+".invalid", util.ErrorNotAllowedFor("decimal", "product qty"))
			}
		}

		poiID, e := common.Decrypt(v.ID)
		if e != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item"))
		}

		poi, e := repository.ValidPurchaseOrderItem(poiID)
		if e != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item"))
		}

		if e = poi.Product.Read("ID"); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		v.TaxableItem = poi.TaxableItem

		if len(gr) > 0 {
			gri := &model.GoodsReceiptItem{GoodsReceipt: gr[0], PurchaseOrderItem: poi}

			if e = gri.Read("GoodsReceipt", "PurchaseOrderItem"); e != nil {
				o.Failure("good_receipt_item.invalid", util.ErrorInvalidData("good receipt item"))
			}

			if v.OrderQty < gri.ReceiveQty {
				o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("product quantity", "receive quantity"))
			}
		}

		if v.OrderQty <= 0 {
			o.Failure("qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("product quantity", "0"))
		}

		if v.UnitPrice < 0 {
			o.Failure("unit_price"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("product unit price", "0"))
		}

		unitPriceInput := v.UnitPrice
		oldQty := poi.OrderQty
		v.TaxAmount = math.Round((unitPriceInput * v.TaxPercentage / 100) * v.OrderQty)
		v.UnitPriceTax = math.Round(unitPriceInput * (100 + v.TaxPercentage) / 100)

		isIncludeTax := v.IncludeTax == 1
		isNotTaxableItem := v.TaxableItem != 1

		if isIncludeTax {
			unitPriceNonTax := math.Round(unitPriceInput * 100 / (100 + v.TaxPercentage))
			unitPriceTax := unitPriceInput

			v.TaxAmount = math.Round((unitPriceTax - unitPriceNonTax) * v.OrderQty)
			v.UnitPriceTax = unitPriceTax
			v.UnitPrice = unitPriceNonTax
		}

		if isNotTaxableItem {
			v.TaxAmount = 0
			v.UnitPriceTax = 0
		}

		poi.OrderQty = v.OrderQty
		poi.UnitPrice = v.UnitPrice
		poi.IncludeTax = v.IncludeTax
		poi.UnitPriceTax = v.UnitPriceTax
		poi.TaxAmount = v.TaxAmount
		poi.Subtotal = v.OrderQty * v.UnitPrice
		poi.Weight = v.OrderQty * v.Product.UnitWeight

		if poi.PurchasePlanItem != nil {
			v.PurchasePlanItem, e = repository.ValidPurchasePlanItem(poi.PurchasePlanItem.ID)
			if e != nil {
				o.Failure("purchase_plan_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase plan item"))
			}

			totalPurchaseQty := v.OrderQty + v.PurchasePlanItem.PurchaseQty - oldQty
			if totalPurchaseQty > v.PurchasePlanItem.PurchasePlanQty {
				o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorEqualLess("total purchase qty", "purchase plan qty"))
			}

			poi.UnitPrice = unitPriceInput
			poi.TaxAmount = 0
			poi.UnitPriceTax = 0
			poi.IncludeTax = 2
			poi.Subtotal = v.OrderQty * poi.UnitPrice
			v.PurchasePlanItem.PurchaseQty -= oldQty
			poi.PurchasePlanItem.PurchaseQty = v.PurchasePlanItem.PurchaseQty
			c.PurchasePlan.TotalPurchaseQty -= oldQty
		}

		// Summarize all the item tax amount
		c.TaxAmount += poi.TaxAmount
		c.TotalPrice = c.TotalPrice + poi.Subtotal
		c.TotalWeight = c.TotalWeight + (poi.OrderQty * v.Product.UnitWeight)

		c.PurchaseOrderItems = append(c.PurchaseOrderItems, poi)

		if val, ok := c.MapDebitNoteItem[v.Product.ID]; ok {
			var unitPrice float64
			if v.TaxAmount == 0 {
				unitPrice = v.UnitPrice
			} else {
				unitPrice = v.UnitPriceTax
			}
			c.TotalPriceDebitNote += common.Rounder(unitPrice*val.ReturnQty, 0.5, 2)
		}

	}

	if len(c.Images) > 4 {
		o.Failure("purchase_order_image.invalid", util.ErrorEqualLess("photo", "4 photos"))
	}

	if len(c.Images) > 0 {
		for _, v := range c.Images {
			if v.ID != "" {
				imageID, e := common.Decrypt(v.ID)
				if e != nil {
					o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item"))
				}

				_, e = repository.ValidPurchaseOrderImage(imageID)
				if e != nil {
					o.Failure("purchase_order_image_id.invalid", util.ErrorInvalidData("purchase order image"))
				}
			}
		}
	}

	c.TotalCharge = c.TotalPrice + c.DeliveryFee + (c.TaxPct * c.TotalPrice / 100) + c.TaxAmount

	return o
}

func (c *updateProductRequest) Messages() map[string]string {
	return map[string]string{
		"purchase_order_items.required": util.ErrorInputRequired("purchase order items"),
	}
}
