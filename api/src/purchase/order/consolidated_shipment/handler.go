// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consolidated_shipment

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/:id", h.detail, auth.AuthorizedFieldPurchaserMobile())
	r.POST("/signature", h.sign, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/print/:id", h.printCount, auth.AuthorizedFieldPurchaserMobile())
	r.POST("", h.consolidate, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/:id", h.update, auth.AuthorizedFieldPurchaserMobile())
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.ConsolidatedShipment
	var total int64
	var warehouseID int64

	warehouse := ctx.QueryParam("warehouse")

	if warehouse != "" {
		if warehouseID, e = common.Decrypt(warehouse); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetConsolidatedShipments(rq, warehouseID); e != nil {
		ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetConsolidatedShipment("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// sign : function to sign consolidated shipment
func (h *Handler) sign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r signRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Sign(r)

	return ctx.Serve(e)
}

// printCount : function to count  copy of print
func (h *Handler) printCount(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r printRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Print(r)

	return ctx.Serve(e)
}

// consolidate : function to consolidate shipment
func (h *Handler) consolidate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r consolidateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Consolidate(r)

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
