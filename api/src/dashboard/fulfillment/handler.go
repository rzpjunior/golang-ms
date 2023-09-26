// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for dashboard overview.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/summary", h.summary, auth.Authorized("dash_fulfillment"))
	r.GET("/product", h.product, auth.Authorized("dash_fulfillment"))
	r.GET("/report", h.report, auth.Authorized("rpt_ful"))
	r.GET("/reset", h.resetCache, auth.Authorized("clr_cache_fulfillment"))
}

// summary : get fulfillment summary
func (h *Handler) summary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	param := ctx.QueryParams()
	date := param.Get("date")
	dateArr := strings.Split(date, "|")
	warehouseID, e := common.Decrypt(param.Get("warehouse"))

	if e == nil {
		if data, e := repository.GetFulfillmentSummary(dateArr, warehouseID); e == nil {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// product : get fulfillment product
func (h *Handler) product(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	param := ctx.QueryParams()
	date := param.Get("date")
	dateArr := strings.Split(date, "|")
	warehouseID, e := common.Decrypt(param.Get("warehouse"))

	if e == nil {
		if data, e := repository.GetUnfulfillmentProduct(dateArr, warehouseID); e == nil {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// report : handler to create report of dashboard fulfillment
func (h *Handler) report(c echo.Context) (e error) {
	var file string

	ctx := c.(*cuxs.Context)
	param := ctx.QueryParams()
	year := param.Get("year")
	warehouseID, e := common.Decrypt(param.Get("warehouse"))
	isExport := param.Get("export") == "1"

	warehouse, e := repository.ValidWarehouse(warehouseID)

	if data, lastUpdatedAt, e := repository.GetReportFulfillment(year, warehouseID); e == nil {
		if isExport {
			if file, e = PrintReportFulfillmentXls(data, warehouse, year, lastUpdatedAt); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// resetCache : handler to reset cache
func (h *Handler) resetCache(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	key := "fulfill"
	if ctx.QueryParams().Get("key") != "" {
		key = ctx.QueryParams().Get("key")
	}

	e = dbredis.Redis.DeleteCacheWhereLike("*" + key + "*")

	return ctx.Serve(e)
}
