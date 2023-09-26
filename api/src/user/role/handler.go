// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized("rol_rdl"))
	r.GET("/filter", h.getFiltered, auth.Authorized("filter_rdl"))
	r.GET("/permission", h.permission, auth.Authorized("filter_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("rol_rdd"))
	r.POST("", h.create, auth.Authorized("rol_crt"))
	r.PUT("/:id", h.update, auth.Authorized("rol_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("rol_arc"))
	r.PUT("/unarchive/:id", h.unArchive, auth.Authorized("rol_urc"))
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetRoles(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) getFiltered(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetFilterRole(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) permission(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetRolePermissions(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetRole("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Save(r)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.ID, e = ctx.Decrypt("id"); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = r.Update()
			}
		}
	}
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

// unArchive : function to set status of archived data into active
func (h *Handler) unArchive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UnArchive(r)
			}
		}
	}

	return ctx.Serve(e)
}
