// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fridge

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/sold_product", h.soldProduct)
}

// soldProduct : function to get sold Product Fridge report
func (h *Handler) soldProduct(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	//area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	soldDateStr := ctx.QueryParam("sold_date")
	soldDateArr := strings.Split(soldDateStr, "|")

	cond := make(map[string]interface{})

	// if area != 0 {
	// 	cond["so.area_id = "] = area
	// }

	if warehouse != 0 {
		cond["bf.warehouse_id = "] = warehouse
	}

	if soldDateStr != "" {
		cond["bf.last_seen_at between "] = soldDateArr
	}

	data, e := getAllProductFridge(cond)
	if e == nil {
		if isExport {
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			var file string
			if file, e = getAllProductXls(soldDateStr, data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}
