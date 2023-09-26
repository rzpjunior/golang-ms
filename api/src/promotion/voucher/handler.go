// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package voucher

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
	r.GET("", h.read, auth.Authorized("vou_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("vou_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("vou_crt"))
	r.POST("/apply", h.apply, auth.Authorized("filter_rdl"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("vou_arc"))
	r.POST("/bulky", h.bulky, auth.Authorized("vou_blk_imp"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Voucher
	var total int64

	if data, total, e = repository.GetVouchers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get requested data based on parameters
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		ctx.ResponseData, e = repository.GetVoucher("id", id)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Voucher
	var total int64

	if data, total, e = repository.GetFilterVoucher(rq); e == nil {
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

// archive : function to set status of active data into archive
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Archive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// apply : function to apply voucher
func (h *Handler) apply(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r applyRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = r.Voucher, nil
	}

	return ctx.Serve(e)
}

// bulky : function to create new data based on input
func (h *Handler) bulky(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r bulkyRequest
	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData = CreateBulky(r)
		}
	}

	return ctx.Serve(e)
}
