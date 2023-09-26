// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse_coverage

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
	r.GET("", h.read, auth.Authorized("wrh_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("wrh_cvrg_crt"))
	r.PUT("/delete/:id", h.delete, auth.Authorized("wrh_cvrg_del"))
	r.PUT("/main_warehouse/:id", h.updateMain, auth.Authorized("wrh_cvrg_upd"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.WarehouseCoverage
	var total int64

	if data, total, e = repository.GetWarehouseCoverages(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.WarehouseCoverage
	var total int64

	if data, total, e = repository.GetFilterWarehouseCoverages(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Save(r)

	return ctx.Serve(e)
}

// delete : function to delete data based on input
func (h *Handler) delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r deleteRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Delete(r)

	return ctx.Serve(e)
}

// updateMain : function to update data based on input
func (h *Handler) updateMain(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateMainRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = UpdateMain(r)

	return ctx.Serve(e)
}
