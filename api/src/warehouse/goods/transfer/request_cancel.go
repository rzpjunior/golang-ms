// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer

import (
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type cancelRequest struct {
	ID   int64  `json:"-"`
	Note string `json:"note" valid:"required"`

	GoodsTransfer     *model.GoodsTransfer       `json:"-"`
	GoodsReceipt      *model.GoodsReceipt        `json:"-"`
	GoodsTransferItem []*model.GoodsTransferItem `json:"-"`
	Stock             []*model.Stock             `json:"-"`
	StockLog          []*model.StockLog          `json:"-"`
	WasteLog          []*model.WasteLog          `json:"-"`

	WarehouseOrigin      *model.Warehouse            `json:"-"`
	WarehouseDestination *model.Warehouse            `json:"-"`
	StockType *model.Glossary   `json:"-"`
	Session   *auth.SessionData `json:"-"`
}

func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var filter, exclude map[string]interface{}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if r.GoodsTransfer, err = repository.ValidGoodsTransfer(r.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("id"))
	}
	if r.GoodsTransfer.Status != 1 {
		o.Failure("status.inactive", util.ErrorDocStatus("goods receipt", "active"))
	} else {
		r.GoodsTransfer.Origin.Read("ID")
	}

	r.StockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", r.GoodsTransfer.StockType)
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	filter = map[string]interface{}{"status": 1, "warehouse_id": r.GoodsTransfer.Origin.ID, "stock_type": r.StockType.ValueInt}
	if _, countStockOpname, e := repository.CheckStockOpnameData(filter, exclude); e == nil && countStockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active ", "stock opname", r.GoodsTransfer.Origin.Code+"-"+r.GoodsTransfer.Origin.Name))
	}

	orSelect.LoadRelated(r.GoodsTransfer, "GoodsTransferItem")

	if r.WarehouseOrigin, err = repository.ValidWarehouse(r.GoodsTransfer.Origin.ID); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse origin"))
		return o
	}

	if r.WarehouseDestination, err = repository.ValidWarehouse(r.GoodsTransfer.Destination.ID); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse destination"))
		return o
	}

	// region wh restriction validation
	warehouseRestriction := make(map[int64]bool)
	if r.Session.Staff.Warehouse.ID != r.WarehouseOrigin.ID {
		if _, err := orSelect.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(r.Session.Staff.WarehouseAccessStr, ",")).All(&r.Session.Staff.WarehouseAccess); err != nil {
			o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse of user"))
		}

		for _,v := range r.Session.Staff.WarehouseAccess{
			warehouseRestriction[v.ID] = true
		}

		if ok,_ := warehouseRestriction[r.WarehouseOrigin.ID];!ok{
			o.Failure("warehouse.invalid", util.ErrorMustBeSame("warehouse of user", "warehouse origin"))
		}
	}
	// endregion

	if r.StockType.ValueName == "good stock" {

		for i, v := range r.GoodsTransfer.GoodsTransferItem {

			r.GoodsTransferItem = append(r.GoodsTransferItem, v)

			stockLog := &model.StockLog{
				Warehouse: r.GoodsTransfer.Origin,
				Product:   v.Product,
				Ref:       r.GoodsTransfer.ID,
				Status:    1,
			}

			if err = stockLog.Read("Warehouse", "Product", "Ref", "Status"); err == nil {
				r.StockLog = append(r.StockLog, stockLog)
			} else {
				o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("stock log"))
			}

			stock := &model.Stock{
				Warehouse: r.GoodsTransfer.Origin,
				Product:   v.Product,
			}

			if err = stock.Read("Warehouse", "Product"); err == nil {
				r.Stock = append(r.Stock, stock)
			} else {
				o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("stock"))
			}
		}
	} else if r.StockType.ValueName == "waste stock" {

		for i, v := range r.GoodsTransfer.GoodsTransferItem {

			r.GoodsTransferItem = append(r.GoodsTransferItem, v)

			wasteLog := &model.WasteLog{
				Warehouse: r.GoodsTransfer.Origin,
				Product:   v.Product,
				Ref:       r.GoodsTransfer.ID,
				Status:    1,
			}

			if err = wasteLog.Read("Warehouse", "Product", "Ref", "Status"); err == nil {
				r.WasteLog = append(r.WasteLog, wasteLog)
			} else {
				o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("stock log"))
			}

			stock := &model.Stock{
				Warehouse: r.GoodsTransfer.Origin,
				Product:   v.Product,
			}

			if err = stock.Read("Warehouse", "Product"); err == nil {
				r.Stock = append(r.Stock, stock)
			} else {
				o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("stock"))
			}
		}
	} else {
		o.Failure("stock_name.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	// region get Goods Receipt
	orSelect.Raw("select * from goods_receipt gr where gr.goods_transfer_id = ? and gr.status = 1", r.GoodsTransfer).QueryRow(&r.GoodsReceipt)
	// endregion

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
