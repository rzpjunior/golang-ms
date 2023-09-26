// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for dashboard overview.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/overview", h.getOverview, auth.Authorized("sls_dash_rd"))
	r.GET("/graph", h.getGraph, auth.Authorized("sls_dash_rd"))
	r.GET("/graph/filter", h.getGraph, auth.Authorized("sls_dash_fil"))
}

func (h *Handler) getOverview(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	param := ctx.QueryParams()
	getDate := param.Get("date")

	data, total, e := GetOverviewByQuery(rq, getDate)

	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) getGraph(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	param := ctx.QueryParams()
	getDate := param.Get("date")

	data, total, e := GetGraphByQuery(rq, getDate)

	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
