// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("are_plc_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("filter_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.PUT("/:id", h.update, auth.Authorized("are_plc_upd"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.AreaPolicy
	var total int64

	if data, total, e = repository.GetAreaPolicies(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.AreaPolicy
	var total int64

	if data, total, e = repository.GetFilterAreaPolicies(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// update : function to update data
func (h *Handler) update(c echo.Context) (e error) {
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

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetAreaPolicy("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
