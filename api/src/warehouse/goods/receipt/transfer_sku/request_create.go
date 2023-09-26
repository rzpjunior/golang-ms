// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer_sku

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type createRequest struct {
	Code            string            `json:"code"`
	GoodReceiptsID  string            `json:"good_receipts_id" valid:"required"`
	WarehouseID     string            `json:"warehouse_id" valid:"required"`
	PurchaseOrderID string            `json:"purchase_order_id"`
	GoodsTransferID string            `json:"goods_transfer_id"`
	InboundType     string            `json:"inbound_type"`
	Products        []*productRequest `json:"products" valid:"required"`

	GoodsReceipt      *model.GoodsReceipt  `json:"-"`
	Warehouse         *model.Warehouse     `json:"-"`
	PurchaseOrder     *model.PurchaseOrder `json:"-"`
	GoodsTransfer     *model.GoodsTransfer `json:"-"`
	TransferSku       *model.TransferSku   `json:"-"`
	TotalTransferQty  float64              `json:"-"`
	TotalWasteQty     float64              `json:"-"`
	RecognitionDateAt time.Time            `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// parent sku
type productRequest struct {
	ID                  string               `json:"id" valid:"required"`
	GoodsReceiptItemQty float64              `json:"gri_qty" valid:"required, numeric"`
	TransferTo          []*transferToRequest `json:"transfer_to" valid:"required"`

	ConvertWeightParent float64        `json:"-"`
	ConvertWeightChild  float64        `json:"-"`
	Discrepancy         float64        `json:"-"`
	Product             *model.Product `json:"-"`
}

// child sku
type transferToRequest struct {
	ProductID   string  `json:"product_id" valid:"required"`
	TransferQty float64 `json:"transfer_qty" valid:"required, numeric"`
	WasteQty    float64 `json:"waste_qty" valid:"required, numeric"`
	WasteReason int8    `json:"waste_reason"`

	Product *model.Product `json:"-"`
	Stock   *model.Stock   `json:"-"`
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var isProductExist = make(map[string]bool)
	var isTransferProductExist = make(map[string]bool)
	var totalTransferQty float64 // for document level
	var totalWasteQty float64    // for document level
	layout := "2006-01-02"
	var err error
	var o1 = orm.NewOrm()
	o1.Using("read_only")

	// region basic validation
	if c.Code, err = util.CheckTable("transfer_sku"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if c.RecognitionDateAt, err = time.Parse(layout, time.Now().Format(layout)); err != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("transfer sku recognition date"))
	}

	if c.GoodReceiptsID != "" {
		goodReceiptsID, err := common.Decrypt(c.GoodReceiptsID)

		if err != nil {
			o.Failure("good_receipts_id.invalid", util.ErrorInvalidData("good_receipts"))
			return o
		}

		if c.GoodsReceipt, err = repository.ValidGoodsReceipt(goodReceiptsID); err != nil {
			o.Failure("good_receipts_id.invalid", util.ErrorInvalidData("good_receipts"))
			return o
		}

		if c.GoodsReceipt.Status != 2 {
			o.Failure("good_receipts_id.invalid", util.ErrorDocStatus("good_receipts", "finished"))
			return o
		}

		stockType, err := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", c.GoodsReceipt.StockType)
		if err != nil {
			o.Failure("stock_type.invalid", util.ErrorInvalidData("stock type"))
		}

		if stockType.ValueName == "waste stock" {
			o.Failure("stock_type.invalid", util.ErrorDocStatus("stock_type", "good stock"))
		}

	}

	if c.WarehouseID != "" {
		warehouseID, err := common.Decrypt(c.WarehouseID)

		if err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}

		if c.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}

		if c.Warehouse.Status != 1 {
			o.Failure("warehouse_id.invalid", util.ErrorActive("warehouse"))
		}

	}

	switch c.InboundType {
	case "purchase_order":
		if c.PurchaseOrderID != "" {
			purchaseOrderID, err := common.Decrypt(c.PurchaseOrderID)

			if err != nil {
				o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase_order"))
			}

			if c.PurchaseOrder, err = repository.ValidPurchaseOrder(purchaseOrderID); err != nil {
				o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase_order"))
			}

			if c.PurchaseOrder.Status == 3 {
				o.Failure("purchase_order_id.invalid", util.ErrorDocStatus("purchase_order_id", "not canceled"))
			}

			c.PurchaseOrder.Supplier.Read("ID")

			if c.GoodsReceipt.PurchaseOrder.ID != c.PurchaseOrder.ID {
				o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order doesn't match goods receipt"))
			}

			var isExistSupplierReturn int8
			o1.Raw("select EXISTS("+
				"select po.id "+
				"from supplier_return sr "+
				"join goods_receipt gr on gr.id = sr.goods_receipt_id "+
				"join purchase_order po on po.id = gr.purchase_order_id "+
				"where po.id = ? and sr.status in (1,2)) ts_sku", c.PurchaseOrder.ID).QueryRow(&isExistSupplierReturn)

			if isExistSupplierReturn != 0 {
				o.Failure("purchase_order_id.invalid", util.ErrorCreateDocStatus("transfer sku", "there is supplier return where status", "active/finished"))
			}
		}
	case "goods_transfer":
		var goodsTransferID int64
		if goodsTransferID, err = common.Decrypt(c.GoodsTransferID); err != nil {
			o.Failure("goods_transfer_id.invalid", util.ErrorInvalidData("goods transfer"))
		}

		if c.GoodsTransfer, err = repository.ValidGoodsTransfer(goodsTransferID); err != nil {
			o.Failure("goods_transfer_id.invalid", util.ErrorInvalidData("goods transfer"))
		}
	}

	if c.GoodsReceipt.Warehouse.ID != c.Warehouse.ID {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse doesn't match goods receipt"))
	}
	// endregion

	for i, v := range c.Products {
		if _, productExist := isProductExist[v.ID]; productExist {
			o.Failure("products_"+strconv.Itoa(i)+"_id.duplicate", util.ErrorDuplicate("product"))
			return o
		}

		var productID int64
		if productID, err = common.Decrypt(v.ID); err != nil {
			o.Failure("products_"+strconv.Itoa(i)+"_id.invalid", util.ErrorInvalidData("product"))
			return o
		}

		if v.Product, err = repository.ValidProduct(productID); err != nil {
			o.Failure("products_"+strconv.Itoa(i)+"_id.invalid", util.ErrorInvalidData("product"))
			return o
		}

		var isExist int8
		o1.Raw("select exists(select ts.id "+
			"from transfer_sku ts "+
			"join transfer_sku_item tsi on ts.id = tsi.transfer_sku_id "+
			"where ts.goods_receipt_id = ? and ts.status in(1,2) and tsi.product_id = ?)", c.GoodsReceipt.ID, v.Product.ID).QueryRow(&isExist)
		if isExist > 0 {
			o.Failure("products_"+strconv.Itoa(i)+"_id.invalid", util.ErrorCreateDoc("transfer sku", "active product"))
			return o
		}

		isProductExist[v.ID] = true
		var initTransferProductExist = make(map[string]bool)
		isTransferProductExist = initTransferProductExist

		// region convert qty to kg weight
		v.ConvertWeightParent = v.GoodsReceiptItemQty * v.Product.UnitWeight
		// endregion

		for j, k := range v.TransferTo {
			if _, transferProductExist := isTransferProductExist[k.ProductID]; transferProductExist {
				o.Failure("products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".duplicate", util.ErrorDuplicate("transfer product"))
				return o
			}

			var transferProductID int64
			if transferProductID, err = common.Decrypt(k.ProductID); err != nil {
				o.Failure("products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".product_id.invalid", util.ErrorInvalidData("transfer product"))
				return o
			}

			if k.Product, err = repository.ValidProduct(transferProductID); err != nil {
				o.Failure("products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+"product_id.invalid", util.ErrorInvalidData("transfer product"))
				return o
			}

			isTransferProductExist[k.ProductID] = true

			if k.WasteQty > 0 {
				v.ConvertWeightChild += k.WasteQty * k.Product.UnitWeight
				if k.WasteReason == 0 {
					o.Failure("waste_reason_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("waste reason"))
				} else {
					// check waste reason in glossary
					_, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "waste_reason", "value_int", k.WasteReason)
					if e != nil {
						o.Failure("waste_reason.invalid", util.ErrorInvalidData("waste_reason"))
					}
				}
			} else {
				k.WasteReason = 0
			}

			v.ConvertWeightChild += k.TransferQty * k.Product.UnitWeight

			// region to add total transfer qty for document level
			if v.Product.ID != k.Product.ID {
				totalTransferQty += k.TransferQty
			}
			// region

			// region to add total waste qty for document level
			if v.Product.ID == k.Product.ID {
				totalWasteQty += k.WasteQty
			}
			// region
		}
		v.Discrepancy = v.ConvertWeightParent - v.ConvertWeightChild
		if v.Discrepancy < -0.001 {
			o.Failure("discrepancy.invalid", util.ErrorGreater("discrepancy", "0"))
		}
	}
	c.TotalTransferQty = totalTransferQty
	c.TotalWasteQty = totalWasteQty

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":      util.ErrorInputRequired("warehouse"),
		"good_receipts_id.required":  util.ErrorInputRequired("good_receipts"),
		"purchase_order_id.required": util.ErrorInputRequired("purchase_order"),
	}

	for i, v := range c.Products {
		messages["products_"+strconv.Itoa(i)+"_id.id.required"] = util.ErrorInputRequired("product")

		for j := range v.TransferTo {
			messages["products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".product_id.required"] = util.ErrorInputRequired("transfer product")
			messages["products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".transfer_qty.required"] = util.ErrorInputRequired("transfer qty")
			messages["products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".waste_qty.required"] = util.ErrorInputRequired("waste qty")
		}
	}

	return messages
}
