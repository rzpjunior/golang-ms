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

// createRequest : struct to hold price set request data
type createRequest struct {
	Code                   string    `json:"-"`
	RequestDateStr         string    `json:"request_date" valid:"required"`
	AreaOriginID           string    `json:"area_origin_id" valid:"required"`
	WarehouseOriginID      string    `json:"warehouse_origin_id" valid:"required"`
	AreaDestinationID      string    `json:"area_destination_id" valid:"required"`
	WarehouseDestinationID string    `json:"warehouse_destination_id" valid:"required"`
	AdditionalCost         float64   `json:"additional_cost"`
	AdditionalCostNote     string    `json:"additional_cost_note"`
	Note                   string    `json:"note"`
	StockTypeID            int8      `json:"stock_type"`
	RequestDate            time.Time `json:"-"`

	AreaOrigin           *model.Area                 `json:"-"`
	WarehouseOrigin      *model.Warehouse            `json:"-"`
	AreaDestination      *model.Area                 `json:"-"`
	WarehouseDestination *model.Warehouse            `json:"-"`
	GoodsTransferItems   []*goodsTransferItemRequest `json:"items" valid:"required"`
	Stocks               []*model.Stock              `json:"-"`
	StockType            *model.Glossary             `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type goodsTransferItemRequest struct {
	GoodsTransferItemID string  `json:"goods_transfer_item_id"`
	ProductID           string  `json:"product_id"`
	TransferQty         float64 `json:"transfer_qty"`
	RequestQty          float64 `json:"request_qty"`
	ReceiveQty          float64 `json:"receive_qty"`
	ReceiveNote         string  `json:"receive_note"`
	UnitCost            float64 `json:"unit_cost"`
	Note                string  `json:"note"`

	GoodsTransferItem *model.GoodsTransferItem `json:"-"`
	Product           *model.Product           `json:"-"`
	Stock             *model.Stock             `json:"-"`
}

// Validate : function to validate uom request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var isProductExist = make(map[string]bool)
	var filter, exclude map[string]interface{}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	layout := "2006-01-02"
	if r.RequestDate, err = time.Parse(layout, r.RequestDateStr); err != nil {
		o.Failure("request_date.invalid", util.ErrorInvalidData("order date"))
	}

	if areaOriginID, err := common.Decrypt(r.AreaOriginID); err == nil {
		if r.AreaOrigin, err = repository.ValidArea(areaOriginID); err != nil {
			o.Failure("area_origin_id.invalid", util.ErrorInvalidData("area origin"))
		}
	} else {
		o.Failure("area_origin_id.invalid", util.ErrorInvalidData("area origin"))
	}

	r.StockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", r.StockTypeID)
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	if warehouseOriginID, err := common.Decrypt(r.WarehouseOriginID); err == nil {
		if r.WarehouseOrigin, err = repository.ValidWarehouse(warehouseOriginID); err == nil {
			filter = map[string]interface{}{"warehouse_id": r.WarehouseOrigin.ID, "status": 1, "stock_type": r.StockType.ValueInt}
			if _, countStockOpname, err := repository.CheckStockOpnameData(filter, exclude); err == nil && countStockOpname > 0 {
				o.Failure("warehouse_origin_id.invalid", util.ErrorRelated("active ", "stock opname", r.WarehouseOrigin.Code+"-"+r.WarehouseOrigin.Name))
			}
		} else {
			o.Failure("warehouse_origin_id.invalid", util.ErrorInvalidData("warehouse origin"))
		}
	} else {
		o.Failure("warehouse_origin_id.invalid", util.ErrorInvalidData("warehouse origin"))
	}

	if areaDestinationID, err := common.Decrypt(r.AreaDestinationID); err == nil {
		if r.AreaDestination, err = repository.ValidArea(areaDestinationID); err != nil {
			o.Failure("area_destination_id.invalid", util.ErrorInvalidData("area destination"))
		}
	} else {
		o.Failure("area_destination_id.invalid", util.ErrorInvalidData("area destination"))
	}

	if warehouseDestinationID, err := common.Decrypt(r.WarehouseDestinationID); err == nil {
		if r.WarehouseDestination, err = repository.ValidWarehouse(warehouseDestinationID); err != nil {
			o.Failure("warehouse_destination_id.invalid", util.ErrorInvalidData("warehouse destination"))
		}
	} else {
		o.Failure("warehouse_destination_id.invalid", util.ErrorInvalidData("warehouse destination"))
	}

	if r.AdditionalCost < 0 {
		o.Failure("additional_cost.greaterequal", util.ErrorEqualGreater("additional cost", "0"))
	}

	if r.AdditionalCost > 0 && r.AdditionalCostNote == "" {
		o.Failure("additional_cost_note.required", util.ErrorInputRequired("additional note"))
	}

	// region wh restriction validation
	warehouseRestriction := make(map[int64]bool)
	if r.Session.Staff.Warehouse.ID != r.WarehouseDestination.ID {
		if _, err := o1.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(r.Session.Staff.WarehouseAccessStr, ",")).All(&r.Session.Staff.WarehouseAccess); err != nil {
			o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse of user"))
		}

		for _,v := range r.Session.Staff.WarehouseAccess{
			warehouseRestriction[v.ID] = true
		}

		if ok,_ := warehouseRestriction[r.WarehouseDestination.ID];!ok{
			o.Failure("warehouse.invalid", util.ErrorMustBeSame("warehouse of user", "warehouse destination"))
		}
	}
	// endregion

	if r.WarehouseOrigin.ID == r.WarehouseDestination.ID{
		o.Failure("warehouse.invalid", util.ErrorInputCannotBeSame("warehouse origin","warehouse destination"))
	}

	for i, v := range r.GoodsTransferItems {
		var productID int64
		if _, exist := isProductExist[v.ProductID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
		}

		if productID, err = common.Decrypt(v.ProductID); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}

		if v.Product, err = repository.ValidProduct(productID); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}
		isProductExist[v.ProductID] = true

		if err = v.Product.Uom.Read("ID"); err != nil {
			o.Failure("uom_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("uom"))
		}
		if v.Product.Uom.DecimalEnabled == 2 {
			if v.RequestQty != float64((int64(v.RequestQty))) {
				o.Failure("request_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("request quantity"))
			}
		}

		if v.RequestQty <= 0 {
			o.Failure("request_qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("request quantity", "0"))
		}

		filter = map[string]interface{}{"product_id": v.Product.ID, "warehouse_id": r.WarehouseOrigin.ID, "status": 1}
		var countStock int64
		var stock []*model.Stock
		if stock, countStock, err = repository.CheckStockData(filter, exclude); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}

		if countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		} else {
			v.Stock = stock[0]
		}

		filter = map[string]interface{}{"product_id": v.Product.ID, "warehouse_id": r.WarehouseDestination.ID, "status": 1}
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}

	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"request_date.required":             util.ErrorInputRequired("request date"),
		"area_origin_id.required":           util.ErrorInputRequired("area origin"),
		"warehouse_origin_id.required":      util.ErrorInputRequired("warehouse origin"),
		"area_destination_id.required":      util.ErrorInputRequired("area destination"),
		"warehouse_destination_id.required": util.ErrorInputRequired("warehouse destination"),
		"eta_date.required":                 util.ErrorInputRequired("eta date"),
		"eta_time.required":                 util.ErrorInputRequired("eta time"),
	}

	return messages
}
