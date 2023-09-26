// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field_purchaser

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/purchase_plan", h.getSummaryPurchasePlan, auth.AuthorizedFieldPurchaserMobile())
}

func (h *Handler) getSummaryPurchasePlan(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	staff, _ := common.Decrypt(ctx.QueryParam("staff"))

	cond := make(map[string]interface{})

	if staff != 0 && staff != 1 {
		cond["pp.assigned_to = "] = staff
	}

	if warehouse != 0 && warehouse != 21 {
		cond["pp.warehouse_id = "] = warehouse
	}

	data, total, e := GetPurchasePlanSummary(cond)
	if e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}
