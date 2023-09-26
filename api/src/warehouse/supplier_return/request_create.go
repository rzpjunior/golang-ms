// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_return

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

// createRequest : struct to hold supplier return request data
type createRequest struct {
	Code string `json:"-"`

	WarehouseID        string                `json:"warehouse_id" valid:"required"`
	SupplierID         string                `json:"supplier_id" valid:"required"`
	GoodReceiptID      string                `json:"good_receipt_id" valid:"required"`
	RecognitionDate    string                `json:"recognition_date" valid:"required"`
	Note               string                `json:"note"`
	TotalPrice         float64               `json:"-"`
	SupplierReturnItem []*SupplierReturnItem `json:"supplier_return_item" valid:"required"`

	SupplierReturn    *model.SupplierReturn
	RecognitionDateAt time.Time
	Warehouse         *model.Warehouse    `json:"-"`
	Supplier          *model.Supplier     `json:"-"`
	GoodsReceipt      *model.GoodsReceipt `json:"-"`

	PurchaseInvoice *model.PurchaseInvoice `json:"-"`
	Session         *auth.SessionData      `json:"-"`
}

type SupplierReturnItem struct {
	ProductID     string  `json:"product_id"`
	ReturnGoodQty float64 `json:"return_good_qty"`
	ReceivedQty   float64 `json:"received_qty"`
	Note          string  `json:"note"`

	Product          *model.Product          `json:"-"`
	GoodsReceiptItem *model.GoodsReceiptItem `json:"-"`
}

// Validate : function to validate supplier return request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q := orm.NewOrm()
	q.Using("read_only")

	var e error
	var whID, supplierID, goodReceiptID int64
	var filter, exclude map[string]interface{}
	var duplicated = make(map[int64]bool)

	layout := "2006-01-02"
	if r.RecognitionDateAt, e = time.Parse(layout, r.RecognitionDate); e != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("recognition date"))
		return o
	}

	// region warehouse definition
	if whID, e = common.Decrypt(r.WarehouseID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		return o
	}

	if r.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		return o
	}
	// endregion

	// region supplier definition
	if supplierID, e = common.Decrypt(r.SupplierID); e != nil {
		o.Failure("supplier_id.invalid", util.ErrorInvalidData("supplier"))
		return o
	}

	if r.Supplier, e = repository.ValidSupplier(supplierID); e != nil {
		o.Failure("supplier_id.invalid", util.ErrorInvalidData("supplier"))
		return o
	}
	// endregion

	// region good receipt definition
	if goodReceiptID, e = common.Decrypt(r.GoodReceiptID); e != nil {
		o.Failure("good_receipt_id.invalid", util.ErrorInvalidData("good receipt"))
		return o
	}

	if r.GoodsReceipt, e = repository.ValidGoodsReceipt(goodReceiptID); e != nil {
		o.Failure("good_receipt_id.invalid", util.ErrorInvalidData("good receipt"))
		return o
	}
	// endregion

	if isExist := q.QueryTable(new(model.SupplierReturn)).Filter("goods_receipt_id", r.GoodsReceipt.ID).Filter("status__in", 1, 2).Exist(); isExist {
		o.Failure("supplier_return_id.invalid", util.ErrorCreateDoc("supplier return", "good receipt"))
		return o
	}

	if r.GoodsReceipt.Status != 2 {
		o.Failure("good_receipt_id.invalid", util.ErrorInvalidData("good receipt"))
		return o
	}

	stockType, err := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "good stock")
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	if isExist := q.QueryTable(new(model.StockOpname)).Filter("warehouse_id", r.Warehouse.ID).Filter("status", 1).Filter("stock_type", stockType.ValueInt).Exist(); isExist {
		o.Failure("stock_opname.active", util.ErrorOneActiveInWarehouse())
		return o
	}

	for i, v := range r.SupplierReturnItem {
		// region product definition
		var pID int64
		if pID, e = common.Decrypt(v.ProductID); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
			return o
		}

		if v.Product, e = repository.ValidProduct(pID); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
			return o
		}
		// endregion

		if duplicated[v.Product.ID] {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorDuplicate("product"))
		}
		if isProductExist := q.QueryTable(new(model.GoodsReceiptItem)).Filter("goods_receipt_id", r.GoodsReceipt.ID).Filter("product_id", v.Product.ID).Exist(); !isProductExist {
			o.Failure("good_receipt_item.invalid", util.ErrorMustExistWarehouse("product", "selected good receipt document"))
			return o
		}

		filter = map[string]interface{}{"product_id": v.Product.ID, "warehouse_id": r.Warehouse.ID, "status": 1, "product__status": 1}
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}

		if v.ReturnGoodQty <= 0 {
			o.Failure("return_good_qty"+strconv.Itoa(i)+".invalid", util.ErrorGreater("return good quantity", "0"))
		}

		if v.ReceivedQty < v.ReturnGoodQty {
			o.Failure("return_good_qty"+strconv.Itoa(i)+".invalid", util.ErrorEqualGreater("received qty", "return good qty"))
		}

		// region get purchase order item data
		v.GoodsReceiptItem = new(model.GoodsReceiptItem)
		if e = q.QueryTable(new(model.GoodsReceiptItem)).Filter("goods_receipt_id", r.GoodsReceipt.ID).Filter("product_id", v.Product.ID).One(v.GoodsReceiptItem); e != nil {
			o.Failure("goods_receipt_item_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("good receipt item"))
		}
		if e = v.GoodsReceiptItem.PurchaseOrderItem.Read("ID"); e != nil {
			o.Failure("purchase_order_item_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase order item"))
		}
		// endregion

		duplicated[v.Product.ID] = true

	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":         util.ErrorInputRequired("warehouse"),
		"supplier_id.required":          util.ErrorInputRequired("supplier"),
		"recognition_date.required":     util.ErrorInputRequired("recognition date"),
		"supplier_return_item.required": util.ErrorInputRequired("supplier return item"),
	}

	return messages
}
