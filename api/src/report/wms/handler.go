// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wms

import (
	"encoding/json"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/stock-log", h.stockLog, auth.Authorized("wrh_rpt_6_dl"))
	r.GET("/waste-log", h.wasteLog, auth.Authorized("wrh_rpt_7_dl"))
	r.GET("/stock", h.stock, auth.Authorized("wrh_rpt_5_dl"))
	r.GET("/goods-receipt-item", h.goodsReceiptItem, auth.Authorized("wrh_rpt_2_dl"))
	r.GET("/delivery-return-item", h.deliveryReturnItem, auth.Authorized("wrh_rpt_8_dl"))
	r.GET("/products", h.products, auth.Authorized("src_rpt_4_dl"))
	r.GET("/delivery-order", h.getDeliveryOrder, auth.Authorized("lgs_rpt_1_dl"))
	r.GET("/item-recap", h.itemRecap, auth.Authorized("src_rpt_1_dl"))
	r.GET("/movement-stock", h.movementStocks, auth.Authorized("wrh_rpt_9_dl"))
	r.GET("/picking", h.picking, auth.Authorized("wrh_rpt_10_dl"))
	r.GET("/picking-order-item", h.pickingOrderItem, auth.Authorized("wrh_rpt_11_dl"))
	r.GET("/goods-transfer-item", h.goodsTransferItem, auth.Authorized("wrh_rpt_12_dl"))
	r.GET("/picking-routing", h.pickingRouting, auth.Authorized("wrh_rpt_14_dl"))
	r.GET("/transfer-sku-item", h.transferSkuItem, auth.Authorized("wrh_rpt_13_dl"))
	r.GET("/packing-recommendation", h.getPackingReport, auth.Authorized("wrh_rpt_1_dl"))
}

// stockLog : function to get stock log report
func (h *Handler) stockLog(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	createdAtStr := ctx.QueryParam("created_at")
	createdAtArr := strings.Split(createdAtStr, "|")
	typeStr := ctx.QueryParam("type")
	refType := ctx.QueryParam("ref_type")
	product, _ := common.Decrypt(ctx.QueryParam("product"))

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["w.id = "] = warehouse
	}

	if createdAtStr != "" {
		cond["date(sl.created_at) "] = createdAtArr
	}

	if typeStr != "" {
		cond["sl.type = "] = typeStr
	}

	if refType != "" {
		cond["sl.ref_type = "] = refType
	}

	if product != 0 {
		cond["sl.product_id = "] = product
	}

	data, e := getStockLog(cond)

	if e == nil {
		if isExport {
			var file string
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			if file, e = GetStockLogXls(data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// wasteLog : function to get waste log report
func (h *Handler) wasteLog(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	createdAtStr := ctx.QueryParam("created_at")
	createdAtArr := strings.Split(createdAtStr, "|")
	typeStr := ctx.QueryParam("type")
	refType := ctx.QueryParam("ref_type")
	product, _ := common.Decrypt(ctx.QueryParam("product"))

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["w.id = "] = warehouse
	}

	if createdAtStr != "" {
		cond["date(wl.created_at) "] = createdAtArr
	}

	if typeStr != "" {
		cond["wl.type = "] = typeStr
	}

	if refType != "" {
		cond["wl.ref_type = "] = refType
	}

	if product != 0 {
		cond["wl.product_id = "] = product
	}

	data, e := getWasteLog(cond)

	if e == nil {
		if isExport {
			var file string
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			if file, e = GetWasteLogXls(data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// stock : function to get stock report
func (h *Handler) stock(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["w.id = "] = warehouse
	}

	data, e := getStock(cond)

	if e != nil {
		return ctx.Serve(e)
	}

	var file string
	if isExport {
		mWarehouse, e := repository.ValidWarehouse(warehouse)
		if file, e = GetStockXls(data, mWarehouse); e != nil {
			return ctx.Serve(e)
		}
	} else {
		ctx.Data(data)
	}

	ctx.Files(file)

	return ctx.Serve(e)
}

// goodsReceiptItem : function to get goods receipt item report
func (h *Handler) goodsReceiptItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))
	ataDateStr := ctx.QueryParam("ata_date")
	ataDateArr := strings.Split(ataDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["gr.warehouse_id = "] = warehouse
	}

	if supplier != 0 {
		cond["po.supplier_id = "] = supplier
	}

	if ataDateStr != "" {
		cond["gr.ata_date "] = ataDateArr
	}

	data, e := getGoodsReceiptItem(cond)

	if e != nil {
		return ctx.Serve(e)
	}

	var file string
	if isExport {
		mWarehouse, e := repository.ValidWarehouse(warehouse)
		if file, e = GetGoodsReceiptItemXls(data, mWarehouse); e != nil {
			return ctx.Serve(e)
		}
	} else {
		ctx.Data(data)
	}

	ctx.Files(file)

	return ctx.Serve(e)
}

// deliveryReturnItem : function to get delivery return item report
func (h *Handler) deliveryReturnItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["w.id = "] = warehouse
	}

	if recognitionDateStr != "" {
		cond["dr.recognition_date "] = recognitionDateArr
	}

	data, e := getDeliveryReturnItem(cond)

	if e == nil {
		if isExport {
			var file string
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			if file, e = GetDeliveryReturnItemXls(data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// products : function to get products report
func (h *Handler) products(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	statusStr := ctx.QueryParam("status")
	salabilityStr, _ := common.Decrypt(ctx.QueryParam("salability"))
	purchasabilityStr, _ := common.Decrypt(ctx.QueryParam("purchasability"))
	storabilityStr, _ := common.Decrypt(ctx.QueryParam("storability"))

	cond := make(map[string]interface{})

	if statusStr != "" && statusStr != "999" {
		cond["p.status = "] = statusStr
	}

	if salabilityStr != 0 {
		cond["p.warehouse_sal "] = salabilityStr
	}

	if purchasabilityStr != 0 {
		cond["p.warehouse_pur "] = purchasabilityStr
	}

	if storabilityStr != 0 {
		cond["p.warehouse_sto "] = storabilityStr
	}

	data, e := getProducts(cond)

	if e == nil {
		if isExport {
			var file string
			if file, e = GetProductsXls(data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// DeliveryOrder : function to get delivery order report
func (h *Handler) getDeliveryOrder(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time
	var warehouse *model.Warehouse
	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	param := ctx.QueryParams()
	wrtID, _ := common.Decrypt(param.Get("wrt_id"))
	warehouseID, _ := common.Decrypt(param.Get("warehouse_id"))
	areaID, _ := common.Decrypt(param.Get("area_id"))

	deliveryDateStr := ctx.QueryParam("delivery_dates")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")

	cond := make(map[string]interface{})

	if deliveryDateStr != "" {
		cond["so.delivery_date between "] = deliveryDateArr
	}

	if wrtID != 0 {
		cond["so.wrt_id = "] = wrtID
	}
	if warehouseID != 0 {
		cond["so.warehouse_id = "] = warehouseID
		warehouse, e = repository.ValidWarehouse(warehouseID)
	}
	if areaID != 0 {
		cond["so.area_id = "] = areaID
	}

	data, e := getDeliveryOrder(rq, cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getDeliveryOrderXls(backdate, data, warehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// itemRecap : function to get item recap report
func (h *Handler) itemRecap(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	param := ctx.QueryParams()
	areaID, _ := common.Decrypt(param.Get("area"))
	warehouseID, _ := common.Decrypt(param.Get("warehouse"))

	area, _ := repository.GetArea("id", areaID)

	deliveryDateStr := ctx.QueryParam("delivery_date")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")

	cond := make(map[string]interface{})

	if deliveryDateStr != "" {
		cond["so.delivery_date "] = deliveryDateArr
	}

	if areaID != 0 {
		cond["a.id = "] = areaID
	}
	if warehouseID != 0 {
		cond["w.id = "] = warehouseID
	}

	data, e := getItemRecap(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = GetItemRecapXls(data, area); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// movementStocks : function to get movement stock report
func (h *Handler) movementStocks(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	param := ctx.QueryParams()
	warehouseID, _ := common.Decrypt(param.Get("warehouse"))

	warehouse, _ := repository.GetWarehouse("id", warehouseID)

	deliveryDateStr := ctx.QueryParam("recognition_date")

	cond := make(map[string]interface{})

	if deliveryDateStr != "" {
		cond["recognition_date"] = deliveryDateStr
	}

	if warehouseID != 0 {
		cond["w.id"] = warehouseID
	}

	data, e := getMovementStock(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = GetMovementStockXls(data, warehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// picking : function to get picking report
func (h *Handler) picking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 {
		cond["so.warehouse_id = "] = warehouse
	}

	if recognitionDateStr != "" {
		cond["so.delivery_date between "] = recognitionDateArr
	}

	data, e := getPicking(cond)
	if e == nil {
		if isExport {
			var file string
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			if file, e = GetPickingXls(data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// picking : function to get picking report
func (h *Handler) pickingOrderItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 {
		cond["po.warehouse_id = "] = warehouse
	}

	if recognitionDateStr != "" {
		cond["po.recognition_date between "] = recognitionDateArr
	}

	data, e := getPickingOrderItem(cond)
	if e != nil {
		return ctx.Serve(e)
	}
	if isExport {
		var file string
		mWarehouse, e := repository.ValidWarehouse(warehouse)
		if file, e = GetPickingOrderItemXls(data, mWarehouse); e == nil {
			ctx.Files(file)
		}
	} else {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// goodsTransferItem : function to get goods transfer item report
func (h *Handler) goodsTransferItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouseOrigin, e := common.Decrypt(ctx.QueryParam("warehouse_origin"))
	if e != nil {
		return ctx.Serve(e)
	}
	warehouseDestination, e := common.Decrypt(ctx.QueryParam("warehouse_destination"))
	if e != nil {
		return ctx.Serve(e)
	}
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")

	cond := make(map[string]interface{})

	if warehouseOrigin != 0 {
		cond["gt.origin_id = "] = warehouseOrigin
	}

	if warehouseDestination != 0 {
		cond["gt.destination_id = "] = warehouseDestination
	}

	if recognitionDateStr != "" {
		cond["gt.recognition_date BETWEEN "] = recognitionDateArr
	}

	data, e := getGoodsTransferItem(cond)
	if e != nil {
		return ctx.Serve(e)
	}
	if isExport {
		var file string
		mWarehouse, e := repository.ValidWarehouse(warehouseOrigin)
		if e != nil {
			return ctx.Serve(e)
		}
		if file, e = GetGoodsTransferItemXls(data, mWarehouse); e != nil {
			return ctx.Serve(e)
		}
		ctx.Files(file)
	} else {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// pickingRouting : function to get picking routing report
func (h *Handler) pickingRouting(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, e := common.Decrypt(ctx.QueryParam("warehouse"))
	if e != nil {
		return ctx.Serve(e)
	}
	product, _ := common.Decrypt(ctx.QueryParam("product"))
	staff, _ := common.Decrypt(ctx.QueryParam("staff"))
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 {
		cond["po.warehouse_id = "] = warehouse
	}

	if recognitionDateStr != "" {
		cond["po.recognition_date between "] = recognitionDateArr
	}

	if product != 0 {
		cond["p.id = "] = product
	}

	if staff != 0 {
		cond["prs.staff_id = "] = staff
	}

	data, e := getPickingRoutingReport(cond)
	if e == nil {
		if isExport {
			var file string
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			if file, e = GetPickingRoutingXls(data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// transferSkuItem : function to get transfer sku item report
func (h *Handler) transferSkuItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	ataDateStr := ctx.QueryParam("recognition_date")
	ataDateArr := strings.Split(ataDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["ts.warehouse_id = "] = warehouse
	}

	if ataDateStr != "" {
		cond["ts.recognition_date "] = ataDateArr
	}

	data, e := getTransferSkuItem(cond)

	if e != nil {
		return ctx.Serve(e)
	}

	var file string
	if isExport {
		mWarehouse, e := repository.ValidWarehouse(warehouse)
		if file, e = GetTransferSkuItemXls(data, mWarehouse); e != nil {
			return ctx.Serve(e)
		}
	} else {
		ctx.Data(data)
	}

	ctx.Files(file)

	return ctx.Serve(e)
}

func (h *Handler) getPackingReport(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var wh *model.Warehouse

	isExport := ctx.QueryParam("export") == "1"
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	deliveryDateStr := ctx.QueryParam("delivery_date")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 {
		cond["po.warehouse_id = "] = warehouse
	}

	if deliveryDateStr != "" {
		cond["po.delivery_date "] = deliveryDateArr
	}

	data, e := getPackingRecommendation(cond)

	wh, _ = repository.ValidWarehouse(warehouse)

	var poID []int64

	for _, v := range data {
		poID = append(poID, v.ID)
	}

	m := mongodb.NewMongo()

	var ps []*model.BarcodeModel
	if len(poID) != 0 {
		filter := bson.D{
			{"packing_order_id", bson.M{"$in": poID}},
		}

		var res []byte
		if res, e = m.GetAllDataWithFilter("Packing_Barcode", filter); e != nil {
			return ctx.Serve(e)
		}

		// region convert byte data to json data
		if e = json.Unmarshal(res, &ps); e != nil {
			return ctx.Serve(e)
		}
		// endregion
	}

	if e == nil {
		if isExport {
			var file string
			if file, e = GetPackingRecommendationXls(ps, wh); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	m.DisconnectMongoClient()
	return ctx.Serve(e)
}
