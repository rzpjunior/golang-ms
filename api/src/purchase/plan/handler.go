// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("ppl_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("ppl_rdd"))
	r.POST("", h.create, auth.Authorized("ppl_crt"))
	r.PUT("/assign/:id", h.assign, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/:id", h.update, auth.Authorized("ppl_upd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("ppl_can"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("ppl_cnf"))
	r.GET("/field_purchaser", h.readInFieldPurchaserApps, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/assign/cancel/:id", h.cancelAssignment, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchasePlan
	var total int64
	var warehouseID int64

	warehouse := ctx.QueryParam("warehouse")

	if warehouse != "" {
		if warehouseID, e = common.Decrypt(warehouse); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetPurchasePlans(rq, warehouseID); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

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

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetPurchasePlan("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// assign purchase order to field purchaser
func (h *Handler) assign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r assignRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Assign(r)

	return ctx.Serve(e)
}

// update : function to update requested data based on parameters
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Update(r)

	return ctx.Serve(e)
}

// cancel : function to cancel purchase plan
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Cancel(r)

	return ctx.Serve(e)
}

// confirm : function to confirm purchase plan
func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Confirm(r)

	return ctx.Serve(e)
}

// readInFieldPurchaserApps : function to get requested data based on parameters
func (h *Handler) readInFieldPurchaserApps(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchasePlan
	var total int64
	var warehouseID int64

	warehouse := ctx.QueryParam("warehouse")

	if warehouse != "" {
		if warehouseID, e = common.Decrypt(warehouse); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetPurchasePlansInFieldPurchaserApps(rq, warehouseID); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// cancelAssignment : function to cancel purchase plan assignment
func (h *Handler) cancelAssignment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelAssignmentRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = CancelAssignment(r)

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchasePlan
	var total int64

	if data, total, e = repository.GetFilterPurchasePlans(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
