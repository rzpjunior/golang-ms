// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package area

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized("are_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("are_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Area
	var total int64

	if data, total, e = repository.GetAreas(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetArea("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	areaId := ctx.QueryParam("area")

	data, total, e := repository.GetFilterArea(ctx.RequestQuery(), areaId)
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
