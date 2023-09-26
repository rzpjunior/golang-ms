// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer

import (
	"git.edenfarm.id/cuxs/orm"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequest : struct to hold price set request data
type updateRequest struct {
	ID                 int64     `json:"-"`
	RecognitionDateStr string    `json:"recognition_date" valid:"required"`
	RequestDateStr     string    `json:"request_date" valid:"required"`
	EtaDateStr         string    `json:"eta_date"`
	EtaTimeStr         string    `json:"eta_time"`
	AdditionalCost     float64   `json:"additional_cost"`
	AdditionalCostNote string    `json:"additional_cost_note"`
	Note               string    `json:"note"`
	RecognitionDate    time.Time `json:"-"`
	RequestDate        time.Time `json:"-"`
	EtaDate            time.Time `json:"-"`
	EtaTime            time.Time `json:"-"`
	TotalWeight        float64   `json:"-"`
	TotalCost          float64   `json:"-"`
	TotalCharge        float64   `json:"-"`
	UpdateAt           time.Time `json:"update_at"`

	AreaOrigin           *model.Area                 `json:"-"`
	WarehouseOrigin      *model.Warehouse            `json:"-"`
	AreaDestination      *model.Area                 `json:"-"`
	WarehouseDestination *model.Warehouse            `json:"-"`
	GoodsTransfer        *model.GoodsTransfer        `json:"-"`
	GoodsTransferItems   []*goodsTransferItemRequest `json:"items" valid:"required"`
	Stock                []*model.Stock              `json:"-"`
	StockType            *model.Glossary             `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var isProductExist = make(map[string]bool)
	var filter, exclude map[string]interface{}
	var gr []*model.GoodsReceipt
	o1 := orm.NewOrm()
	o1.Using("read_only")

	if r.GoodsTransfer, err = repository.ValidGoodsTransfer(r.ID); err != nil {
		o.Failure("goods_transfer.invalid", util.ErrorInvalidData("goods transfer"))
		return o
	}

	filter = map[string]interface{}{"purchase_order_id": r.GoodsTransfer.ID, "status__in": 1}
	gr, _, _ = repository.CheckGoodsReceiptData(filter, exclude)

	layout := "2006-01-02"
	if r.RecognitionDate, err = time.Parse(layout, r.RecognitionDateStr); err != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("recognition date"))
	}

	if r.RequestDate, err = time.Parse(layout, r.RequestDateStr); err != nil {
		o.Failure("request_date.invalid", util.ErrorInvalidData("request date"))
	}

	if r.WarehouseOrigin, err = repository.ValidWarehouse(r.GoodsTransfer.Origin.ID); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse origin"))
		return o
	}

	if r.WarehouseDestination, err = repository.ValidWarehouse(r.GoodsTransfer.Destination.ID); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse destination"))
		return o
	}

	warehouseRestriction := make(map[int64]bool)
	if _, err := o1.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(r.Session.Staff.WarehouseAccessStr, ",")).All(&r.Session.Staff.WarehouseAccess); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse of user"))
	}

	for _,v := range r.Session.Staff.WarehouseAccess{
		warehouseRestriction[v.ID] = true
	}
	switch r.GoodsTransfer.Status {
	case 1:
		if r.EtaDate, err = time.Parse(layout, r.EtaDateStr); err != nil {
			o.Failure("eta_date.invalid", util.ErrorInvalidData("eta date"))
		}

		if r.EtaTime, err = time.Parse("15:04", r.EtaTimeStr); err != nil {
			o.Failure("eta_time.invalid", util.ErrorInvalidData("eta time"))
		}
		if r.Session.Staff.Warehouse.ID != r.WarehouseOrigin.ID {
			if ok,_ := warehouseRestriction[r.WarehouseOrigin.ID];!ok{
				o.Failure("warehouse.invalid", util.ErrorMustBeSame("warehouse of user", "warehouse origin"))
			}
		}
	case 5:
		if r.Session.Staff.Warehouse.ID != r.WarehouseDestination.ID {
			if ok,_ := warehouseRestriction[r.WarehouseDestination.ID];!ok{
				o.Failure("warehouse.invalid", util.ErrorMustBeSame("warehouse of user", "warehouse destination"))
			}
		}
	}

	if r.EtaDate.Before(r.RecognitionDate) {
		o.Failure("eta_date.greater", util.ErrorEqualLater("eta date", "departure date"))
	}

	if r.AdditionalCost < 0 {
		o.Failure("additional_cost.greaterequal", util.ErrorEqualGreater("additional cost", "0"))
	}

	if r.AdditionalCost > 0 && r.AdditionalCostNote == "" {
		o.Failure("additional_cost_note.required", util.ErrorInputRequired("additional note"))
	}

	r.StockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", r.GoodsTransfer.StockType)
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	filter = map[string]interface{}{"warehouse_id": r.GoodsTransfer.Origin.ID, "status": 1, "stock_type": r.StockType.ValueInt}
	if _, countStockOpname, err := repository.CheckStockOpnameData(filter, exclude); err == nil && countStockOpname > 0 {
		o.Failure("warehouse_origin_id.invalid", util.ErrorRelated("active ", "stock opname", r.WarehouseOrigin.Code+"-"+r.WarehouseOrigin.Name))
	}

	for i, v := range r.GoodsTransferItems {
		if _, exist := isProductExist[v.ProductID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
		} else {
			var productID int64
			if productID, err = common.Decrypt(v.ProductID); err != nil {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			}
			if v.Product, err = repository.ValidProduct(productID); err != nil {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			}
			isProductExist[v.ProductID] = true

			if v.TransferQty > v.RequestQty {
				o.Failure("request_qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("request quantity", "transfer quantity"))
			}

			if err = v.Product.Uom.Read("ID"); err != nil {
				o.Failure("uom_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("uom"))
			}
			if r.GoodsTransfer.Status == 1 {
				if v.Product.Uom.DecimalEnabled == 2 {
					if v.TransferQty != float64((int64(v.TransferQty))) {
						o.Failure("transfer_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("transfer quantity"))
					}
				}
				if v.TransferQty < 0 {
					o.Failure("transfer_qty"+strconv.Itoa(i)+".greater", util.ErrorEqualGreater("transfer quantity", "0"))
				}
			} else {
				if v.Product.Uom.DecimalEnabled == 2 {
					if v.RequestQty != float64((int64(v.RequestQty))) {
						o.Failure("request_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("request quantity"))
					}
				}
				if v.RequestQty <= 0 {
					o.Failure("request_qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("request quantity", "0"))
				}
			}

			if len(gr) > 0 {
				gri := &model.GoodsReceiptItem{GoodsReceipt: gr[0], Product: v.Product}

				if err = gri.Read("GoodsReceipt", "Product"); err != nil {
					o.Failure("good_receipt_item.invalid", util.ErrorInvalidData("product"))
				}

				if v.TransferQty < gri.ReceiveQty {
					o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("product quantity", "transfer quantity"))
				}
			}

			filter = map[string]interface{}{"warehouse_id": r.GoodsTransfer.Origin.ID, "status": 1}
			if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
				o.Failure("warehouse_origin_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
			}

			filter = map[string]interface{}{"warehouse_id": r.GoodsTransfer.Destination.ID, "status": 1}
			if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
				o.Failure("warehouse_destination_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
			}

			r.TotalWeight = r.TotalWeight + (v.Product.UnitWeight * v.TransferQty)
			r.TotalCost = r.TotalCost + (v.UnitCost * v.TransferQty)

		}
	}

	r.TotalCharge = r.TotalCost + r.AdditionalCost

	return o
}

// Messages : function to return error validation messages
func (r *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"recognition_date.required":         util.ErrorInputRequired("recognition date"),
		"request_date.required":             util.ErrorInputRequired("request date"),
		"area_origin_id.required":           util.ErrorInputRequired("area origin"),
		"warehouse_origin_id.required":      util.ErrorInputRequired("warehouse origin"),
		"area_destination_id.required":      util.ErrorInputRequired("area destination"),
		"warehouse_destination_id.required": util.ErrorInputRequired("warehouse destination"),
	}

	return messages
}
