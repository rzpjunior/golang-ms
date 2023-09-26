// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adm_division

import (
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("adm_division_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.AdmDivision
	var total int64

	if data, total, e = repository.GetAdmDivisions(rq, ctx.QueryParams().Get("groupby")); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.AdmDivision
	var total int64

	// polygon is a param to get the polygon and centroid of the adm division
	polygon := ctx.QueryParam("polygon")

	if data, total, e = repository.GetFilterAdmDivisions(rq, ctx.QueryParams().Get("groupby"), polygon); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
