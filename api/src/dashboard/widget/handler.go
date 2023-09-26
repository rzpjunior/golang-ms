// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package widget

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for dashboard overview.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/so-wrt", h.getSalesOrderWithWRT, auth.Authorized("opr_dash_rd"))
	r.GET("/idle-picker", h.getIdlePicking, auth.Authorized("opr_dash_rd"))
}

func (h *Handler) getSalesOrderWithWRT(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	area, _ := common.Decrypt(ctx.QueryParam("area"))
	wrt, _ := common.Decrypt(ctx.QueryParam("wrt"))
	deliveryDateStr := ctx.QueryParam("delivery_date")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")

	cond := make(map[string]interface{})

	if area != 0 && area != 1 {
		cond["so.area_id = "] = area
	}

	if warehouse != 0 && warehouse != 21 {
		cond["so.warehouse_id = "] = warehouse
	}

	if deliveryDateStr != "" {
		cond["so.delivery_date "] = deliveryDateArr
	}

	if wrt != 0 {
		cond["so.wrt_id = "] = wrt
	}

	typeRequest := ctx.QueryParam("type_request")
	switch typeRequest {
	case "sales_order":
		data, total, e := getGrandTotalSalesOrderWithWRT(cond)
		if e == nil {
			ctx.Data(data, total)
		}
	case "picking_order":
		data, total, e := getTotalPickingOrderWithWRT(cond)
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) getIdlePicking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))

	cond := make(map[string]interface{})

	if warehouse != 0 && warehouse != 21 {
		cond["w.id = "] = warehouse
	}

	data, total, e := getIdlePicking(cond)
	if e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}
