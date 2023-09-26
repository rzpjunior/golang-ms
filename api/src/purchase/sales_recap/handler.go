// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_recap

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("sls_rcp_rdl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	var (
		deliveryDate            time.Time
		warehouseID, categoryID int64
	)

	ctx := c.(*cuxs.Context)
	params := ctx.QueryParams()

	if params.Get("delivery_date") != "" {
		deliveryDate, e = time.Parse("2006-01-02", params.Get("delivery_date"))
	}

	if params.Get("warehouse") != "" {
		warehouseID, e = common.Decrypt(params.Get("warehouse"))
	}

	if params.Get("category") != "" {
		categoryID, e = common.Decrypt(params.Get("category"))
	}

	if e != nil {
		return ctx.Serve(e)
	}

	var data []*model.SalesRecap
	var total int64

	if data, total, e = repository.GetSalesRecaps(deliveryDate, warehouseID, categoryID); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
