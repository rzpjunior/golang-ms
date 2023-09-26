// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/delivery-order-item", h.reportDeliveryOrderItem, auth.Authorized("lgs_rpt_2_dl"))
	r.GET("/packing", h.getPackingReport, auth.Authorized("wrh_rpt_1_dl"))
	r.GET("/pricing-inbound-item", h.pricingInboundItem, auth.Authorized("pri_rpt_1_dl"))
	r.GET("/price-change-history", h.getPriceChangeHistoryReport, auth.Authorized("pri_rpt_2_dl"))
}

func (h *Handler) reportDeliveryOrderItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var backdate time.Time
	var mWarehouse *model.Warehouse
	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	cond := make(map[string]interface{})

	wrt, _ := common.Decrypt(ctx.QueryParam("wrt_id"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse_id"))
	area, _ := common.Decrypt(ctx.QueryParam("area_id"))

	if wrt != 0 {
		cond["so.wrt_id="] = wrt
	}
	if warehouse != 0 {
		cond["so.warehouse_id="] = warehouse
		mWarehouse, e = repository.ValidWarehouse(warehouse)
	}

	if area != 0 {
		cond["so.area_id="] = area
	}

	deliveryDateStr := ctx.QueryParam("recognition_dates")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")

	if deliveryDateStr != "" {
		cond["so.delivery_date between "] = deliveryDateArr
	}

	data, e := getDeliveryOrderItem(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getDeliveryOrderItemXls(backdate, data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

func (h *Handler) getPackingReport(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var backdate time.Time
	var wh *model.Warehouse

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse_id"))
	deliveryDateStr := ctx.QueryParam("delivery_date")

	data, e := getPackingOrderReport(deliveryDateStr, warehouse)

	wh, _ = repository.ValidWarehouse(warehouse)

	if e == nil {
		if isExport {
			var file string
			if file, e = DownloadPackingReportXls(backdate, data, wh); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// pricingInboundItem : function to get Pricing Inbound Item report
func (h *Handler) pricingInboundItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	etaDateStr := ctx.QueryParam("eta_date")
	etaDateArr := strings.Split(etaDateStr, "|")
	orderDateStr := ctx.QueryParam("order_date")
	orderDateArr := strings.Split(orderDateStr, "|")

	cond := make(map[string]interface{})
	cond2 := make(map[string]interface{})

	if area != 0 {
		cond["wd.area_id = "] = area
		cond2["wd.area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if warehouse != 0 {
		cond["gr.warehouse_id = "] = warehouse
		cond2["gr.warehouse_id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if etaDateStr != "" {
		cond["gt.eta_date between "] = etaDateArr
		cond2["po.eta_date between "] = etaDateArr
		note += fmt.Sprintf("Eta Date : %s - %s; ", etaDateArr[0], etaDateArr[1])
	}

	if orderDateStr != "" {
		cond["gt.recognition_date between "] = orderDateArr
		cond2["po.recognition_date between "] = orderDateArr
		note += fmt.Sprintf("Order Date : %s - %s; ", orderDateArr[0], orderDateArr[1])
	}

	// if user did'nt use warehouse filter, exclude PO to warehouse type ECF (7)
	if _, ok := cond2["gr.warehouse_id = "]; !ok {
		cond2["wd.warehouse_type != "] = 7
	}

	data, e := getPricingInboundItem(cond, cond2)
	if e != nil {
		return ctx.Serve(e)
	}

	if isExport {
		var file string
		mArea, e := repository.ValidArea(area)
		if file, e = GetPricingInboundItem(etaDateStr, data, mArea, session.Staff, note); e != nil {
			return ctx.Serve(e)
		}
		ctx.Files(file)
	} else {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

func (h *Handler) getPriceChangeHistoryReport(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var backdate time.Time
	var countPriceSet int
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	cond := make(map[string]interface{})

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	priceSetIDStr := ctx.QueryParam("price_set_id")
	priceSetIDArr := strings.Split(priceSetIDStr, ".")
	countPriceSet = len(priceSetIDArr)

	createdAtStr := ctx.QueryParam("created_at")
	createdAtArr := strings.Split(createdAtStr, "|")

	if priceSetIDStr != "" {
		for i, v := range priceSetIDArr {
			priceSetID, _ := common.Decrypt(v)
			priceSetIDArr[i] = strconv.Itoa(int(priceSetID))
		}
		cond["p.price_set_id IN "] = priceSetIDArr
		var listPriceSet []*model.PriceSet
		filterCond := make(map[string]interface{})
		excludeCond := make(map[string]interface{})
		filterCond["id__in"] = priceSetIDArr
		var listpriceSetName string

		if listPriceSet, _, e = repository.CheckPriceSetData(filterCond, excludeCond); e != nil {
			return ctx.Serve(e)
		}
		for _, v := range listPriceSet {
			listpriceSetName += v.Name + ", "
		}
		note += fmt.Sprintf("Price Set : %s; ", listpriceSetName)
	}

	if createdAtStr != "" {
		note += fmt.Sprintf("Created Date : %s - %s; ", createdAtArr[0], createdAtArr[1])
		createdAtArr[0] = createdAtArr[0] + " 00:00:01"
		createdAtArr[1] = createdAtArr[1] + " 23:59:59"
		cond["pl.created_at BETWEEN "] = createdAtArr
	}

	data, e := getPriceChangeHistoryReport(cond, countPriceSet)

	if e != nil {
		return ctx.Serve(e)
	}

	if isExport {
		var file string
		if file, e = DownloadPriceChangeHistoryXls(backdate, data, session.Staff, note); e == nil {
			ctx.Files(file)
		}
	} else {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}
