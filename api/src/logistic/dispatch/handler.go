// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dispatch

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.getDispatch, auth.Authorized("dsp_rdl"))
	r.PUT("/:id", h.updateCourier, auth.Authorized("dsp_asg_cou"))
	r.PUT("/vendor/:id", h.updateVendor, auth.Authorized("dsp_asg_cou"))
	r.PUT("/scan", h.scanDispatch, auth.Authorized("dsp_scn"))

	r.GET("/courier", h.getCourier, auth.Authorized("cou_vdr_rdl"))
}

// getDispatch : function to get requested data based on parameters
func (h *Handler) getDispatch(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingOrderAssign
	var total int64

	if data, total, e = repository.GetPickingOrderAssignsforDispatch(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// getCourier : function to get requested data based on parameters
func (h *Handler) getCourier(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Courier
	var total int64

	if data, total, e = repository.GetCouriers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// updateCourier : function to unarchive requested data based on parameters
func (h *Handler) updateCourier(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Update(r)
			}
		}
	}

	return ctx.Serve(e)
}

// updateVendor : function to unarchive requested data based on parameters
func (h *Handler) updateVendor(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateVendorRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateVendor(r)
			}
		}
	}

	return ctx.Serve(e)
}

// scanDispatch : function to unarchive requested data based on parameters
func (h *Handler) scanDispatch(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r scanRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = ScanDispatch(r)
		}
	}

	return ctx.Serve(e)
}
